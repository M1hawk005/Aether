import axios from 'axios';
import MockAdapter from 'axios-mock-adapter';
import { AetherClient } from '../src/client';

describe('AetherClient', () => {
  let mock: MockAdapter;
  let client: AetherClient;

  beforeAll(() => {
    mock = new MockAdapter(axios);
  });

  beforeEach(() => {
    mock.reset();
    client = new AetherClient('http://localhost:5000');
  });

  afterAll(() => {
    mock.restore();
  });

  describe('initialization', () => {
    it('uses default values when no args passed', () => {
      const defaultClient = new AetherClient();
      expect(defaultClient).toBeDefined();
    });
  });

  describe('ping', () => {
    it('returns true on 200', async () => {
      mock.onGet('/api/resolve/ping_test').reply(200, {});
      const result = await client.ping();
      expect(result).toBe(true);
    });

    it('returns true on 404', async () => {
      mock.onGet('/api/resolve/ping_test').reply(404, {});
      const result = await client.ping();
      expect(result).toBe(true);
    });

    it('returns false on 500', async () => {
      mock.onGet('/api/resolve/ping_test').reply(500, {});
      const result = await client.ping();
      expect(result).toBe(false);
    });
  });

  describe('publish', () => {
    it('throws if missing args', async () => {
      await expect(client.publish('', '')).rejects.toThrow('Publisher and payload are required');
    });

    it('returns data on success', async () => {
      mock.onPost('/api/publish').reply(200, { ok: true, capsule_id: '123' });
      const result = await client.publish('alice', 'hello');
      expect(result.ok).toBe(true);
    });

    it('throws daemon error if present', async () => {
      mock.onPost('/api/publish').reply(400, { error: 'Bad data' });
      await expect(client.publish('alice', 'hello')).rejects.toThrow('Bad data');
    });

    it('throws network error if no response', async () => {
      mock.onPost('/api/publish').networkError();
      await expect(client.publish('alice', 'hello')).rejects.toThrow('Network Error');
    });
  });

  describe('resolve', () => {
    it('throws if missing target', async () => {
      await expect(client.resolve('')).rejects.toThrow('Target alias is required');
    });

    it('returns data on success', async () => {
      mock.onGet('/api/resolve/alice').reply(200, { ok: true, capsules: [] });
      const result = await client.resolve('alice');
      expect(result.ok).toBe(true);
    });

    it('handles 404 cleanly', async () => {
      mock.onGet('/api/resolve/bob').reply(404, {});
      const result = await client.resolve('bob');
      expect(result.ok).toBe(false);
      expect(result.error).toBe('No capsules found');
    });

    it('throws on 500', async () => {
      mock.onGet('/api/resolve/charlie').reply(500, {});
      await expect(client.resolve('charlie')).rejects.toThrow('Failed to resolve target');
    });
  });

  describe('fetch', () => {
    it('throws if missing capsuleId', async () => {
      await expect(client.fetch('')).rejects.toThrow('Capsule ID is required');
    });

    it('returns data on success', async () => {
      mock.onGet('/api/fetch/123').reply(200, { ok: true, data: 'hello' });
      const result = await client.fetch('123');
      expect(result.ok).toBe(true);
    });

    it('handles 404 cleanly without retrying', async () => {
      mock.onGet('/api/fetch/123').reply(404, {});
      const result = await client.fetch('123');
      expect(result.ok).toBe(false);
      expect(result.error).toContain('timed out');
    });

    it('throws on 500', async () => {
      mock.onGet('/api/fetch/123').reply(500, {});
      await expect(client.fetch('123')).rejects.toThrow('Failed to fetch capsule');
    });
  });

  describe('fetchByUri', () => {
    it('throws if missing uri', async () => {
      await expect(client.fetchByUri('')).rejects.toThrow();
    });

    it('removes aether:// prefix and calls fetch', async () => {
      mock.onGet('/api/fetch/123').reply(200, { ok: true, data: 'hello' });
      const result = await client.fetchByUri('aether://123');
      expect(result.ok).toBe(true);
    });

    it('removes aether: prefix and calls fetch', async () => {
      mock.onGet('/api/fetch/123').reply(200, { ok: true, data: 'hello' });
      const result = await client.fetchByUri('aether:123');
      expect(result.ok).toBe(true);
    });

    it('calls fetch without prefix', async () => {
      mock.onGet('/api/fetch/123').reply(200, { ok: true, data: 'hello' });
      const result = await client.fetchByUri('123');
      expect(result.ok).toBe(true);
    });
  });

  describe('checkAvailability', () => {
    it('throws if missing capsuleId', async () => {
      await expect(client.checkAvailability('')).rejects.toThrow();
    });

    it('returns data on success', async () => {
      mock.onGet('/api/availability/123').reply(200, { ok: true, seeds: 5 });
      const result = await client.checkAvailability('123');
      expect(result.seeds).toBe(5);
    });

    it('throws on error', async () => {
      mock.onGet('/api/availability/123').reply(500, {});
      await expect(client.checkAvailability('123')).rejects.toThrow();
    });
  });

  describe('deleteLocal', () => {
    it('throws if missing capsuleId', async () => {
      await expect(client.deleteLocal('')).rejects.toThrow();
    });

    it('returns success', async () => {
      mock.onPost('/api/delete/123').reply(200, { ok: true });
      const result = await client.deleteLocal('123');
      expect(result.ok).toBe(true);
    });

    it('handles 404', async () => {
      mock.onPost('/api/delete/123').reply(404, {});
      const result = await client.deleteLocal('123');
      expect(result.ok).toBe(false);
    });

    it('throws on error', async () => {
      mock.onPost('/api/delete/123').reply(500, {});
      await expect(client.deleteLocal('123')).rejects.toThrow();
    });
  });

  describe('revokeGlobal', () => {
    it('throws if missing capsuleId', async () => {
      await expect(client.revokeGlobal('')).rejects.toThrow();
    });

    it('returns success', async () => {
      mock.onPost('/api/revoke/123').reply(200, { ok: true });
      const result = await client.revokeGlobal('123');
      expect(result.ok).toBe(true);
    });

    it('throws on error', async () => {
      mock.onPost('/api/revoke/123').reply(500, {});
      await expect(client.revokeGlobal('123')).rejects.toThrow();
    });
  });

  describe('claimUsername', () => {
    it('throws if missing args', async () => {
      await expect(client.claimUsername('', '')).rejects.toThrow();
    });

    it('returns success', async () => {
      mock.onPost('/api/claim-username').reply(200, { ok: true });
      const result = await client.claimUsername('alice', 'key');
      expect(result.ok).toBe(true);
    });

    it('throws on error', async () => {
      mock.onPost('/api/claim-username').reply(500, {});
      await expect(client.claimUsername('alice', 'key')).rejects.toThrow();
    });
  });

  describe('releaseUsername', () => {
    it('throws if missing args', async () => {
      await expect(client.releaseUsername('', '')).rejects.toThrow();
    });

    it('returns success', async () => {
      mock.onPost('/api/release-username').reply(200, { ok: true });
      const result = await client.releaseUsername('alice', 'key');
      expect(result.ok).toBe(true);
    });

    it('throws on error', async () => {
      mock.onPost('/api/release-username').reply(500, {});
      await expect(client.releaseUsername('alice', 'key')).rejects.toThrow();
    });
  });

  describe('getStatus', () => {
    it('returns status', async () => {
      mock.onGet('/api/status').reply(200, { ok: true, peers: 2 });
      const result = await client.getStatus();
      expect(result.peers).toBe(2);
    });

    it('throws on error', async () => {
      mock.onGet('/api/status').reply(500, {});
      await expect(client.getStatus()).rejects.toThrow();
    });
  });

  describe('getReputation', () => {
    it('returns reputation', async () => {
      mock.onGet('/api/reputation').reply(200, { ok: true, reputation: {} });
      const result = await client.getReputation();
      expect(result.ok).toBe(true);
    });

    it('throws on error', async () => {
      mock.onGet('/api/reputation').reply(500, {});
      await expect(client.getReputation()).rejects.toThrow();
    });
  });
});
