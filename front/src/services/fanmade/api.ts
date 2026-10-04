import {
  createMediaApi,
  uploadFile,
  uploadFiles,
  cancelUpload,
  cancelAllUploads,
  type MediaItem
} from "../media/api/mediaApi.js";

export type FanmadeMediaItem = MediaItem & {
  mediaid?: string | number;
  mediaId?: string | number;
  id?: string | number;
  entityid?: string | number;
  entityId?: string | number;
  entitytype?: string;
  entityType?: string;
  creatorid?: string | number;
  creatorId?: string | number;
  userid?: string | number;
  userId?: string | number;
  captionlang?: string | null;
  captionLang?: string | null;
  createdAt?: string | number | Date | null;
  updatedAt?: string | number | Date | null;
  [key: string]: unknown;
};

const fanmadeApi = createMediaApi("fanmade");

export const fetchMedia = fanmadeApi.fetchMedia.bind(fanmadeApi);
export const deleteMedia = fanmadeApi.deleteMedia.bind(fanmadeApi);
export const postMediaFanmade = fanmadeApi.postMedia.bind(fanmadeApi);

export { uploadFile, uploadFiles, cancelUpload, cancelAllUploads };
export type { MediaItem };
