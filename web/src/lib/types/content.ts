export type Artist = {
	id: string;
	name: string;
	picture_url: string;
	scrobble_count: number | null;
};

export type Album = {
	id: string;
	title: string;
	picture_url: string;
	scrobble_count: number | null;
	artists: Artist[] | null;
};

export type Track = {
	id: string;
	title: string;
	artists: Artist[] | null;
	picture_url: string;
	album: Album;
	scrobble_count: number | null;
	scrobbled_at: Date | null;
};
