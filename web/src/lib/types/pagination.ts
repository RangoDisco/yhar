export type Paginated<T> = {
	results: T[];
	pagination: {
		total_count: number;
		has_next_page: boolean;
	};
};
