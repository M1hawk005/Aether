export interface PublishOptions {
  visibility?: 'public' | 'private' | 'unlisted';
  persistent?: boolean;
}

export interface PublishRequest {
  publisher: string;
  payload: string;
  visibility?: string;
  persistent?: boolean;
}

export interface PublishResponse {
  ok: boolean;
  capsule_id?: string;
  error?: string;
}

export interface CapsuleHeader {
  capsule_id: string;
  created_at: string;
  envelope_hash: string;
  manifest_hash: string;
  owner_key_id: string;
  publisher: string;
  signature: string;
  version: number;
  visibility: string;
}

export interface ResolveResponse {
  ok: boolean;
  username: string;
  capsules: CapsuleHeader[];
  error?: string;
}

export interface FetchResponse {
  ok: boolean;
  data?: string;
  error?: string;
}

export interface AvailabilityResponse {
  ok: boolean;
  capsule_id: string;
  seeds: number;
  error?: string;
}

export interface ClaimRequest {
  username: string;
  key_id: string;
}

export interface StatusResponse {
  ok: boolean;
  peers: number;
  peer_id: string;
}

export interface ReputationResponse {
  ok: boolean;
  reputation: Record<string, number>;
}
