import { connect, Identity, Signer } from '@hyperledger/fabric-gateway';
import { credentials, Client } from '@grpc/grpc-js';
import * as fs from 'fs';
import * as path from 'path';
import { createSign } from 'crypto';

async function newGateway() {
  
  const ccp = JSON.parse(fs.readFileSync('connection-org1.json', 'utf8'));


  const peerCert = fs.readFileSync(path.resolve('crypto/peerOrganizations/org1.DC.com/peers/peer0.org1.DC.com/tls/ca.crt'));
  const grpcCreds = credentials.createSsl(peerCert);
  const grpcClient = new Client('localhost:7051', grpcCreds, {
    'grpc.ssl_target_name_override': 'peer0.org1.DC.com'
  });

  
  const certPem = fs.readFileSync(path.resolve('wallet/appUserCert.pem'));
  const keyPem  = fs.readFileSync(path.resolve('wallet/appUserKey.pem'));
  const userId  = { mspId: 'Org1MSP', credentials: certPem };
  const userSigner = newSigner(keyPem.toString());

  return connect({
    client: grpcClient,       // required
    identity: userId,         // required
    signer: userSigner,       // required
    
  });
}

(async () => {
  const gateway = await newGateway();
  const network = await gateway.getNetwork('mychannel');
  const contract = network.getContract('DC');

  await gateway.close();
})();

export { newGateway };
function newSigner(keyPem: string) {
    return async (digest: Uint8Array): Promise<Uint8Array> => {
        const sign = createSign('sha256');
        sign.update(digest);
        sign.end();
        const signature = sign.sign(keyPem);
        return new Uint8Array(signature);
    };
}
