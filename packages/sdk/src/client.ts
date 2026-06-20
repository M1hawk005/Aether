import axios, { AxiosInstance } from 'axios';
import axiosRetry from 'axios-retry';
import {
  PublishOptions,
  PublishResponse,
  ResolveResponse,
  FetchResponse,
  AvailabilityResponse,
  StatusResponse,
  ReputationResponse
} from './types';

export class AetherClient {
  private api: AxiosInstance;

  constructor(daemonUrl: string = 'http://localhost:5000', timeoutMs: number = 10000) {
    this.api = axios.create({
      baseURL: daemonUrl,
      timeout: timeoutMs,
    });

    // Configure automatic retries for 5xx errors and network drops (ECONNREFUSED)
    axiosRetry(this.api, {
      retries: 3,
      retryDelay: axiosRetry.exponentialDelay,
      retryCondition: (error) => {
        // Do NOT retry 404/408 (e.g. P2P fetch timeouts)
        if (error.response && [404, 408].includes(error.response.status)) {
          return false;
        }
        // Retry on 5xx or network errors (like connection refused)
        return Boolean(axiosRetry.isNetworkOrIdempotentRequestError(error) || 
               (error.response && error.response.status >= 500));
      }
    });
  }

  async ping(): Promise<boolean> {
    try {
      await this.api.get('/api/resolve/ping_test');
      return true;
    } catch (error: any) {
      if (error.response && error.response.status === 404) {
        return true;
      }
      return false;
    }
  }

  async publish(publisher: string, payload: string, options: PublishOptions = {}): Promise<PublishResponse> {
    if (!publisher || !payload) throw new Error("Publisher and payload are required");
    try {
      const res = await this.api.post<PublishResponse>('/api/publish', {
        publisher,
        payload,
        visibility: options.visibility || 'public',
        persistent: options.persistent ?? true,
      });
      return res.data;
    } catch (error: any) {
      if (error.response?.data?.error) {
        throw new Error(error.response.data.error);
      }
      throw error;
    }
  }

  async resolve(targetAlias: string): Promise<ResolveResponse> {
    if (!targetAlias) throw new Error("Target alias is required");
    try {
      const res = await this.api.get<ResolveResponse>(`/api/resolve/${encodeURIComponent(targetAlias)}`);
      return res.data;
    } catch (error: any) {
      if (error.response?.status === 404) {
        return { ok: false, username: targetAlias, capsules: [], error: 'No capsules found' };
      }
      throw new Error(`Failed to resolve target: ${error.message}`);
    }
  }

  async fetch(capsuleId: string): Promise<FetchResponse> {
    if (!capsuleId) throw new Error("Capsule ID is required");
    try {
      const res = await this.api.get<FetchResponse>(`/api/fetch/${encodeURIComponent(capsuleId)}`);
      return res.data;
    } catch (error: any) {
      if (error.response?.status === 404) {
        return { ok: false, error: 'Not cached locally and P2P fetch timed out' };
      }
      throw new Error(`Failed to fetch capsule: ${error.message}`);
    }
  }

  async fetchByUri(uri: string): Promise<FetchResponse> {
    if (!uri) throw new Error("URI is required");
    let capsuleId = uri;
    if (capsuleId.startsWith("aether://")) {
      capsuleId = capsuleId.slice("aether://".length);
    } else if (capsuleId.startsWith("aether:")) {
      capsuleId = capsuleId.slice("aether:".length);
    }
    return this.fetch(capsuleId);
  }

  async checkAvailability(capsuleId: string): Promise<AvailabilityResponse> {
    if (!capsuleId) throw new Error("Capsule ID is required");
    try {
      const res = await this.api.get<AvailabilityResponse>(`/api/availability/${encodeURIComponent(capsuleId)}`);
      return res.data;
    } catch (error: any) {
      throw new Error(`Failed to check availability: ${error.message}`);
    }
  }

  async deleteLocal(capsuleId: string): Promise<{ ok: boolean; message?: string; error?: string }> {
    if (!capsuleId) throw new Error("Capsule ID is required");
    try {
      const res = await this.api.post(`/api/delete/${encodeURIComponent(capsuleId)}`);
      return res.data;
    } catch (error: any) {
      if (error.response?.status === 404) {
        return { ok: false, error: 'Not found locally' };
      }
      throw new Error(`Failed to delete local capsule: ${error.message}`);
    }
  }

  async revokeGlobal(capsuleId: string): Promise<{ ok: boolean; message?: string; error?: string }> {
    if (!capsuleId) throw new Error("Capsule ID is required");
    try {
      const res = await this.api.post(`/api/revoke/${encodeURIComponent(capsuleId)}`);
      return res.data;
    } catch (error: any) {
      throw new Error(`Failed to broadcast revocation: ${error.message}`);
    }
  }

  async claimUsername(username: string, keyId: string): Promise<{ ok: boolean; message?: string; error?: string }> {
    if (!username || !keyId) throw new Error("Username and key ID are required");
    try {
      const res = await this.api.post(`/api/claim-username`, { username, key_id: keyId });
      return res.data;
    } catch (error: any) {
      throw new Error(`Failed to claim username: ${error.message}`);
    }
  }

  async releaseUsername(username: string, keyId: string): Promise<{ ok: boolean; message?: string; error?: string }> {
    if (!username || !keyId) throw new Error("Username and key ID are required");
    try {
      const res = await this.api.post(`/api/release-username`, { username, key_id: keyId });
      return res.data;
    } catch (error: any) {
      throw new Error(`Failed to release username: ${error.message}`);
    }
  }

  async getStatus(): Promise<StatusResponse> {
    try {
      const res = await this.api.get<StatusResponse>(`/api/status`);
      return res.data;
    } catch (error: any) {
      throw new Error(`Failed to get status: ${error.message}`);
    }
  }

  async getReputation(): Promise<ReputationResponse> {
    try {
      const res = await this.api.get<ReputationResponse>(`/api/reputation`);
      return res.data;
    } catch (error: any) {
      throw new Error(`Failed to get reputation: ${error.message}`);
    }
  }
}
