export interface Song {
    songid: string | number;
    songId?: string | number;
    artistid?: string | number;
    artistId?: string | number;
    albumid?: string | number;
    albumId?: string | number;
    title?: string;
    genre?: string;
    duration?: string | number;
    description?: string | null;
    poster?: string | null;
    audioUrl?: string | null;
    audioURL?: string | null;
    audioextn?: string;
    audioExtn?: string;
    published?: boolean;
    plays?: number;
    uploadedAt?: string | number | Date;
    language?: string;
    liked?: boolean;
    _playBtn?: HTMLButtonElement;
    [key: string]: unknown;
}

export interface Playlist {
    playlistid: string | number;
    playlistId?: string | number;
    userid?: string | number;
    userId?: string | number;
    name?: string;
    description?: string;
    songs?: Array<string | Song>;
    createdAt?: string | number | Date;
    updatedAt?: string | number | Date;
    duration?: number;
    isCompilation?: boolean;
    copyrights?: string;
    coverUrl?: string | null;
    coverURL?: string | null;
    [key: string]: unknown;
}

export interface Album {
    albumid: string | number;
    albumId?: string | number;
    artistid?: string | number;
    artistId?: string | number;
    title?: string;
    description?: string;
    published?: boolean;
    releaseDate?: string | number | Date;
    songs?: Array<string | Song>;
    coverUrl?: string | null;
    coverURL?: string | null;
    [key: string]: unknown;
}

export interface Player {
    play: (song: Song, idx: number) => void;
    setQueue?: (songs: Song[]) => void;
}

export type NotifyOptions = {
    type?: "info" | "error" | "success";
};