import { EventEmitter } from 'eventemitter3';
import { WebSocketMonitor } from './WebSocketMonitor';
import { RulesEngine } from '../rules-engine/RulesEngine';
import { MultiChannelDelivery } from '../channels/MultiChannelDelivery';
import { APIServer } from '../api/APIServer';
import { EventProcessor } from './EventProcessor';
import { MetricsCollector } from './MetricsCollector';
import { logger } from '../utils/logger';

export interface NotificationPlatformConfig {
  wsMonitor: WebSocketMonitor;
  rulesEngine: RulesEngine;
  delivery: MultiChannelDelivery;
  apiServer: APIServer;
}

export class NotificationPlatform extends EventEmitter {
  private wsMonitor: WebSocketMonitor;
  private rulesEngine: RulesEngine;
  private delivery: MultiChannelDelivery;
  private apiServer: APIServer;
  private eventProcessor: EventProcessor;
  private metricsCollector: MetricsCollector;
  private isRunning: boolean = false;

  constructor(config: NotificationPlatformConfig) {
    super();
    this.wsMonitor = config.wsMonitor;
    this.rulesEngine = config.rulesEngine;
    this.delivery = config.delivery;
    this.apiServer = config.apiServer;
    this.eventProcessor = new EventProcessor(this.rulesEngine);
    this.metricsCollector = new MetricsCollector();
  }

  async start(): Promise<void> {
    if (this.isRunning) {
      throw new Error('Platform is already running');
    }

    logger.info('Starting notification platform components...');

    // Start WebSocket monitoring
    await this.wsMonitor.connect();
    this.setupEventHandlers();

    // Initialize delivery channels
    await this.delivery.initialize();

    // Start API server
    await this.apiServer.start();

    // Start metrics collection
    this.metricsCollector.start();

    this.isRunning = true;
    this.emit('started');
    
    logger.info('All components started successfully');
  }

  async stop(): Promise<void> {
    if (!this.isRunning) {
      return;
    }

    logger.info('Stopping notification platform...');

    // Stop WebSocket monitoring
    await this.wsMonitor.disconnect();

    // Stop delivery channels
    await this.delivery.shutdown();

    // Stop API server
    await this.apiServer.stop();

    // Stop metrics collection
    this.metricsCollector.stop();

    this.isRunning = false;
    this.emit('stopped');
    
    logger.info('Platform stopped successfully');
  }

  private setupEventHandlers(): void {
    // Handle trades
    this.wsMonitor.on('trade', async (data) => {
      const notifications = await this.eventProcessor.processTrade(data);
      await this.sendNotifications(notifications);
    });

    // Handle liquidations
    this.wsMonitor.on('liquidation', async (data) => {
      const notifications = await this.eventProcessor.processLiquidation(data);
      await this.sendNotifications(notifications);
    });

    // Handle funding
    this.wsMonitor.on('funding', async (data) => {
      const notifications = await this.eventProcessor.processFunding(data);
      await this.sendNotifications(notifications);
    });

    // Handle large orders
    this.wsMonitor.on('largeOrder', async (data) => {
      const notifications = await this.eventProcessor.processLargeOrder(data);
      await this.sendNotifications(notifications);
    });

    // Handle price alerts
    this.wsMonitor.on('priceAlert', async (data) => {
      const notifications = await this.eventProcessor.processPriceAlert(data);
      await this.sendNotifications(notifications);
    });

    // Handle errors
    this.wsMonitor.on('error', (error) => {
      logger.error('WebSocket error:', error);
      this.emit('error', error);
    });

    // Handle reconnection
    this.wsMonitor.on('reconnected', () => {
      logger.info('WebSocket reconnected');
      this.emit('reconnected');
    });
  }

  private async sendNotifications(notifications: any[]): Promise<void> {
    for (const notification of notifications) {
      try {
        await this.delivery.send(notification);
        this.metricsCollector.recordNotification(notification);
      } catch (error) {
        logger.error('Failed to send notification:', error);
        this.metricsCollector.recordError(error);
      }
    }
  }

  public getMetrics(): any {
    return this.metricsCollector.getMetrics();
  }

  public isHealthy(): boolean {
    return this.isRunning && 
           this.wsMonitor.isConnected() && 
           this.delivery.isHealthy();
  }
}