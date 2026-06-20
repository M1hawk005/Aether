import { AetherClient } from './src/index';

async function run() {
  const client = new AetherClient();
  console.log('Testing Aether Node.js SDK...');
  
  const isOnline = await client.ping();
  console.log(`Daemon online: ${isOnline}`);
  
  if (!isOnline) {
    console.error('Please ensure the Aether Gateway is running and unlocked.');
    return;
  }

  // Publisher username
  // NOTE: This identity must be already created and unlocked in the local daemon!
  const publisher = 'neo';

  console.log('\n--- Publishing Payload ---');
  let capsuleId = '';
  try {
    const pubRes = await client.publish(publisher, JSON.stringify({ hello: 'sdk_world', timestamp: Date.now() }));
    console.log('Publish result:', pubRes);
    if (!pubRes.ok || !pubRes.capsule_id) {
        throw new Error('Failed to get capsule_id');
    }
    capsuleId = pubRes.capsule_id;
  } catch (err: any) {
    console.error('Publish error:', err.message);
    if (err.message.includes('not found')) {
      console.error(`Please create identity '${publisher}' in the Gateway UI first.`);
    }
    return;
  }

  console.log('\n--- Resolving Alias ---');
  try {
    const resolveRes = await client.resolve(publisher);
    console.log(`Found ${resolveRes.capsules.length} capsules for ${publisher}`);
  } catch (err: any) {
    console.error('Resolve error:', err.message);
  }

  console.log(`\n--- Fetching Payload for ${capsuleId} ---`);
  try {
    const fetchRes = await client.fetch(capsuleId);
    console.log('Fetch result:', fetchRes);
  } catch (err: any) {
    console.error('Fetch error:', err.message);
  }

  console.log(`\n--- Checking Availability for ${capsuleId} ---`);
  try {
    const availRes = await client.checkAvailability(capsuleId);
    console.log(`Seeds found: ${availRes.seeds}`);
  } catch (err: any) {
    console.error('Availability error:', err.message);
  }
}

run();
