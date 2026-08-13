export type Artist = {
	id: string;
	name: string;
	music_brainz_id: string | null;
	picture_url: string;
	scrobble_count: number | null;
};

export type Album = {
	id: string;
	music_brainz_id: string | null;
	title: string;
	picture_url: string;
	scrobble_count: number | null;
	artists: Artist[] | null;
};

export type Scrobble = {
	id: string;
	scrobbled_at: Date;
	track: Track;
};

export type Track = {
	id: string;
	title: string;
	artists: Artist[] | null;
	picture_url: string;
	album: Album;
	scrobble_count: number;
};
