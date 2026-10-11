import { page } from '$app/state';
import { goto } from '$app/navigation';

export const PER_PAGE = 10;

export const handlePeriodChange = (period: string) => {
	const query = new URLSearchParams(page.url.searchParams.toString());
	query.set('period', period);
	query.set('page', '1');
	goto(`?${query.toString()}`, { reset: false });
};

export const handlePageChange = (newPage: number) => {
	console.log("PAGE CHANGE", newPage);
	const query = new URLSearchParams(page.url.searchParams.toString());
	query.set('page', newPage.toString());
	goto(`?${query.toString()}`, { reset: false });
};
