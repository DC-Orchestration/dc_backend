import express from 'express';
import { startRedisListener } from './redisSubscriber';
import { newGateway } from './fabricGateway';

const app = express();
app.use(express.json());


startRedisListener().catch(console.error);


app.get('/audit/user/:userId', async (req, res) => {
  const gateway = await newGateway();
  const network = await gateway.getNetwork('mychannel');
  const contract = network.getContract('DC');
  const result = await contract.evaluateTransaction('getByUser', req.params.userId);
  await gateway.close();
  res.json(JSON.parse(result.toString()));
});


app.get('/audit/type/:changeType', async (req, res) => {
  const gateway = await newGateway();
  const network = await gateway.getNetwork('mychannel');
  const contract = network.getContract('DC');
  const result = await contract.evaluateTransaction('getByChangeType', req.params.changeType);
  await gateway.close();
  res.json(JSON.parse(result.toString()));
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => console.log(`DC backend listening on port ${PORT}`));
