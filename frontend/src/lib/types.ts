export interface Connection {
  id: string;
  name: string;
  host: string;
  port: number;
  accessKey: string;
  secretKey: string;
  region: string;
  bucket: string;
  useSSL: boolean;
  pathStyle: boolean;
}

export interface S3Object {
  key: string;
  name: string;
  isFolder: boolean;
  size: number;
  lastModified: string;
  storageClass: string;
}

export type TransferType = 'upload' | 'download';
export type TransferStatus = 'pending' | 'in_progress' | 'completed' | 'failed' | 'cancelled';

export interface Transfer {
  id: string;
  fileName: string;
  type: TransferType;
  status: TransferStatus;
  bytesTotal: number;
  bytesDone: number;
  percentage: number;
  error: string;
}

export interface Settings {
  maxParallel: number;
}

export type SortField = 'name' | 'size' | 'lastModified';
export type SortDir = 'asc' | 'desc';

export interface Toast {
  id: string;
  type: 'success' | 'error' | 'info';
  message: string;
}
