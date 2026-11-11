export type ContentType = 'artists' | 'albums' | 'tracks';

export type Artist = {
	id: string;
	name: string;
	music_brainz_id: string | null;
	scrobble_count: number | null;
	picturePath: string | null;
	pictureDomain: string | null;
};

export type Album = {
	id: string;
	music_brainz_id: string | null;
	title: string;
	scrobble_count: number | null;
	artists: Artist[] | null;
	picturePath: string | null;
	pictureDomain: string | null;
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
	album: Album;
	scrobble_count: number;
	picturePath: string | null;
	pictureDomain: string | null;
};

export type TimelineResult = {
	listenedCount: number;
	periodInterval: string;
};

export type TotalScrobbleDurationAndCount = {
	totalScrobbles: number;
	totalDuration: number;
};

export type AvgDailyScrobbles = {
	avgCount: number;
};

export type TotalDistinctCount = {
	totalCount: number;
};

export type SearchResult = {
	contentType: ContentType;
	contentID: string;
	title: string;
	path: string;
	domain: string;
};
