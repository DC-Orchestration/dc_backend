import { createClient } from 'redis';
import { newGateway } from './fabricGateway';

async function startRedisListener() {
  const redis = createClient({ url: 'redis://localhost:6379' });
  await redis.connect();
  await redis.subscribe('doc:changes', async (message: string) => {
    const { docId, userId, changeType, summary, prevHash } = JSON.parse(message);
    const gateway = await newGateway();
    const network = await gateway.getNetwork('mychannel');
    const contract = network.getContract('DC');  // your chaincode name
    await contract.submitTransaction('addAuditLog', docId, userId, changeType, summary, prevHash);
    console.log(`Audit entry for ${docId} by ${userId}`);
    await gateway.close();
  });
}

export { startRedisListener };
