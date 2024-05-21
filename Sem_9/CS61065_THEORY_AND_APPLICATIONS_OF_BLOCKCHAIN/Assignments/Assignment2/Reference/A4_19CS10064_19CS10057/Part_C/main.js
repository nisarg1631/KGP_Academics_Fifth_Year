//  Assignment No. - 4 Part - C
//  Hyperledger Fabric
//  CS61065 - Theory and Applications of Blockchain
//  Semester - 7 (Autumn 2022-23)
//  Group Members - Vanshita Garg (19CS10064) & Shristi Singh (19CS10057)

/**
 * To run this file, directory structuew should be as follows:
 * - Assignment-4
 *      - A4_19CS10064_19CS10057
 *          - Part-A
 *              - p1.go and its dependencies
 *          - Part-B
 *              - p2.go and its dependencies
 *          - Part-C
 *              - main.js
 *              - package.json
 *              - package-lock.json
 *              - node_modules
 *              - wallet
 *      - fabric-samples
 *          - test-network
 * 
 * Note: Please change the paths in this file according to your directory structure.
 */

const FabricCAServices = require('fabric-ca-client');
const { Wallets, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const prompt = require('prompt-sync')({ sigint: true });
const sleep = require('sleep');

async function main() {
    try {
        // load the network configuration for each of the two peers
        // for peer1
        const ccpPath_1 = path.resolve(__dirname, '..', '..', 'fabric-samples', 'test-network', 'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
        const fileExists_1 = fs.existsSync(ccpPath_1);
        if (!fileExists_1) {
            throw new Error(`no such file or directory: ${ccpPath_1}`);
        }
        const ccp_1 = JSON.parse(fs.readFileSync(ccpPath_1, 'utf8'));

        // for peer2
        const ccpPath_2 = path.resolve(__dirname, '..', '..', 'fabric-samples', 'test-network', 'organizations', 'peerOrganizations', 'org2.example.com', 'connection-org2.json');
        const fileExists_2 = fs.existsSync(ccpPath_2);
        console.log('fileExists: ', fileExists_2);
        if (!fileExists_2) {
            throw new Error(`no such file or directory: ${ccpPath_2}`);
        }
        const ccp_2 = JSON.parse(fs.readFileSync(ccpPath_2, 'utf8'));


        // Create a new CA client for interacting with the CA for peer 1.
        const caInfo_1 = ccp_1.certificateAuthorities['ca.org1.example.com'];
        const caTLSCACerts_1 = caInfo_1.tlsCACerts.pem;
        const ca_1 = new FabricCAServices(caInfo_1.url, { trustedRoots: caTLSCACerts_1, verify: false }, caInfo_1.caName);

        // Create a new CA client for interacting with the CA for peer 2.
        const caInfo_2 = ccp_2.certificateAuthorities['ca.org2.example.com'];
        const caTLSCACerts_2 = caInfo_2.tlsCACerts.pem;
        const ca_2 = new FabricCAServices(caInfo_2.url, { trustedRoots: caTLSCACerts_2, verify: false }, caInfo_2.caName);

        // create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user for peer 1.
        identity_1 = await wallet.get('admin_1');
        if (identity_1) {
            console.log('An identity for the user "admin" already exists in the wallet');
        } else {                              
            // enroll the admin user for peer 1, and import the new identity into the wallet.
            const enrollment = await ca_1.enroll({ enrollmentID: 'admin', enrollmentSecret: 'adminpw' });
            const x509Identity = {
                credentials: {
                    certificate: enrollment.certificate,
                    privateKey: enrollment.key.toBytes(),
                },
                mspId: 'Org1MSP',
                type: 'X.509',
            };
            await wallet.put('admin_1', x509Identity);
            console.log('Successfully enrolled user "admin_1" and imported it into the wallet');
            identity_1 = await wallet.get('admin_1');
        }
        // Check to see if we've already enrolled the user for peer 2.
        identity_2 = await wallet.get('admin_2');
        if (identity_2) {
            console.log('An identity for the user "admin" already exists in the wallet');
        } else {                              
            // enroll the admin user for peer 1, and import the new identity into the wallet.
            const enrollment = await ca_2.enroll({ enrollmentID: 'admin', enrollmentSecret: 'adminpw' });
            const x509Identity = {
                credentials: {
                    certificate: enrollment.certificate,
                    privateKey: enrollment.key.toBytes(),
                },
                mspId: 'Org1MSP',
                type: 'X.509',
            };
            await wallet.put('admin_2', x509Identity);
            console.log('Successfully enrolled user "admin_2" and imported it into the wallet');
            identity_2 = await wallet.get('admin_2');
        }

        // user registration for peer 1
        var userIdentity_1 = await wallet.get('appUser_1');
        if (userIdentity_1) {
            console.log('An identity for the user "appUser_1" already exists in the wallet');
        } else {
            // build a user object for authenticating with the CA
            const provider = wallet.getProviderRegistry().getProvider(identity_1.type);
            const adminUser = await provider.getUserContext(identity_1, 'admin');

            // Register the user, enroll the user, and import the new identity into the wallet.
            // check if user is already registered with the CA
            try {
                const secret = await ca_1.register({
                    affiliation: 'org1.department1',
                    enrollmentID: 'appUser_1',
                    role: 'client'
                }, adminUser);
                console.log('Successfully registered user "appUser_1"');
                const enrollment = await ca_1.enroll({
                    enrollmentID: 'appUser_1',
                    enrollmentSecret: secret,
                });
                const x509Identity = {
                    credentials: {
                        certificate: enrollment.certificate,
                        privateKey: enrollment.key.toBytes(),
                    },
                    mspId: 'Org1MSP',
                    type: 'X.509',
                };
                await wallet.put('appUser_1', x509Identity);
            } catch (error) {}
            console.log('Successfully enrolled user "appUser_1" and imported it into the wallet');
            userIdentity_1 = await wallet.get('appUser_1');
        }

        // user registration for peer 2
        var userIdentity_2 = await wallet.get('appUser_2');
        if (userIdentity_2) {
            console.log('An identity for the user "appUser_2" already exists in the wallet');
        } else {
            // build a user object for authenticating with the CA
            const provider = wallet.getProviderRegistry().getProvider(identity_2.type);
            const adminUser = await provider.getUserContext(identity_2, 'admin');

            // Register the user, enroll the user, and import the new identity into the wallet.
            // check if user is already registered with the CA
            try {
                const secret = await ca_2.register({
                    affiliation: 'org2.department1',
                    enrollmentID: 'appUser_2',
                    role: 'client'
                }, adminUser);
                console.log('Successfully registered user "appUser_2"');
                const enrollment = await ca_2.enroll({
                    enrollmentID: 'appUser_2',
                    enrollmentSecret: secret,
                });
                const x509Identity = {
                    credentials: {
                        certificate: enrollment.certificate,
                        privateKey: enrollment.key.toBytes(),
                    },
                    mspId: 'Org2MSP',
                    type: 'X.509',
                };
                await wallet.put('appUser_2', x509Identity);
            } catch (error) {}
            console.log('Successfully enrolled user "appUser_2" and imported it into the wallet');
            userIdentity_2 = await wallet.get('appUser_2');
        }
        console.log("vanshita");
        const gateway_1 = new Gateway();
        await gateway_1.connect(ccp_1, { wallet, identity: 'appUser_1', discovery: { enabled: true, asLocalhost: true } });
        const gateway_2 = new Gateway();
        await gateway_2.connect(ccp_2, { wallet, identity: 'appUser_2', discovery: { enabled: true, asLocalhost: true } });
        // console.log (gateway_1);
        // console.log (gateway_2);

        const network_1 = await gateway_1.getNetwork('mychannel');
        const network_2 = await gateway_2.getNetwork('mychannel');

        // console.log (network_1);
        // console.log (network_2);

        const contract_1 = network_1.getContract('partB');
        const contract_2 = network_2.getContract('partB');

        // handle query alternatively by peer 1 and peer 2
        // maintain a chance to query by peer 1 and peer 2 alternatively
        var queryCount = 0;
        while (true) {
            // take input from user
            var input = prompt('Enter your command: ');
            var inputArray = input.split(" ");
            // check valid input
            if (inputArray.length != 1) {
                console.log('Invalid input');
                continue;
            }
            // convert inputArray[0] to upper case
            inputArray[0] = inputArray[0].toUpperCase();
            if (inputArray[0] == 'INSERT') {
                // take input the number to be inserted in the binary tree
                var number = prompt('Enter the number to be inserted: ');
                // check valid input
                if (isNaN(number)) {
                    console.log('Invalid input');
                    continue;
                }
                var result
                if (queryCount == 0) {
                    result = await contract_1.submitTransaction('Insert', number);
                } else {
                    result = await contract_2.submitTransaction('Insert', number);
                }
            } else if (inputArray[0] == 'DELETE') {
                // take input the number to be deleted from the binary tree
                var number = prompt('Enter the number to be deleted: ');
                // check valid input
                if (isNaN(number)) {
                    console.log('Invalid input');
                    continue;
                }
                var result
                if (queryCount == 0) {
                    result = await contract_1.submitTransaction('Delete', number);
                } else {
                    result = await contract_2.submitTransaction('Delete', number);
                }
            } else if (inputArray[0] == 'INORDER') {
                var result
                if (queryCount == 0) {
                    result = await contract_1.evaluateTransaction('Inorder');
                } else {
                    result = await contract_2.evaluateTransaction('Inorder');
                }
                console.log(result.toString());
            } else if (inputArray[0] == 'PREORDER') {
                var result
                if (queryCount == 0) {
                    result = await contract_1.evaluateTransaction('Preorder');
                } else {
                    result = await contract_2.evaluateTransaction('Preorder');
                }
                console.log(result.toString());
            } else if (inputArray[0] == 'TREEHEIGHT') {
                var result
                if (queryCount == 0) {
                    result = await contract_1.evaluateTransaction('TreeHeight');
                } else {
                    result = await contract_2.evaluateTransaction('TreeHeight');
                }
                console.log(result.toString());
            } else if (inputArray[0] == "EXIT") {
                break;
            } else {
                console.log('Invalid input');
                continue;
            }
            // sleep 2 seconds and output query executed successfully, wait for 2 seconds
            // to make sure the query is executed by peer 1 or peer 2 alternatively
            queryCount = (queryCount + 1) % 2;
            console.log('Query executed successfully, please wait for 2 seconds before executing next query');
            // await sleep(2000);
        }

        await gateway_1.disconnect();
        await gateway_2.disconnect();

    } catch (error) {
        console.error(`Failed to enroll "user" or "admin": ${error}`);
        process.exit(1);
    }
}

main(); 