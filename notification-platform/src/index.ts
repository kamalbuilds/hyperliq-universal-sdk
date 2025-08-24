import { NotificationPlatform } from './core/NotificationPlatform';
import { WebSocketMonitor } from './core/WebSocketMonitor';
import { RulesEngine } from './rules-engine/RulesEngine';
import { MultiChannelDelivery } from './channels/MultiChannelDelivery';
import { APIServer } from './api/APIServer';
import { config } from './config';
import { logger } from './utils/logger';

async function main() {
  try {
    logger.info('Starting Hyperliquid Notification Platform...');

    // Initialize core components
    const wsMonitor = new WebSocketMonitor(config.hyperliquid.wsUrl);
    const rulesEngine = new RulesEngine();
    const delivery = new MultiChannelDelivery();
    const apiServer = new APIServer();

    // Create main platform
    const platform = new NotificationPlatform({
      wsMonitor,
      rulesEngine,
      delivery,
      apiServer
    });

    // Start the platform
    await platform.start();

    logger.info('Notification Platform started successfully');

    // Graceful shutdown
    process.on('SIGTERM', async () => {
      logger.info('SIGTERM received, shutting down gracefully...');
      await platform.stop();
      process.exit(0);
    });

    process.on('SIGINT', async () => {
      logger.info('SIGINT received, shutting down gracefully...');
      await platform.stop();
      process.exit(0);
    });

  } catch (error) {
    logger.error('Failed to start notification platform', error);
    process.exit(1);
  }
}

main().catch(console.error);