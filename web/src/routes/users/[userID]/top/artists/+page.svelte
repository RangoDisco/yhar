<script lang="ts">
    import * as Tabs from "$lib/components/ui/tabs/index";
    import * as Pagination from "$lib/components/ui/pagination/index";
    import ContentListItem from "$lib/components/top/ContentListItem.svelte";
    import {goto} from "$app/navigation";

    let {data} = $props();
    const periods = ["week", "month", "year", "overall"];

    const handlePeriodChange = (period: string) => {
        goto(`?period=${period}&page=1`, {keepFocus: true});
    };

    const handlePageChange = (page: number) => {
        goto(`?period=${data.period}&page=${page}`);
    };

</script>
<Tabs.Root value={data.period} onValueChange={handlePeriodChange}>
    <Tabs.List class="w-full">
        {#each periods as period}
            <Tabs.Trigger value={period}>{period.charAt(0).toUpperCase() + period.slice(1)}</Tabs.Trigger>
        {/each}
    </Tabs.List>
    {#each periods as period}
        <Tabs.Content value={period} class="flex flex-col gap-8">
            {#each data.artists.result as artist, i}
                <ContentListItem index={i} contentID={artist.id} title={artist.name} parents={[]}
                                 pictureUrl={artist.picture_url}
                                 scrobbleCount={artist.scrobble_count}
                                 mode="artists"
                                 contentType="artists"/>
            {/each}
            <Pagination.Root count={data.artists.pagination.total_count} perPage={10} page={data.page}
                             onPageChange={handlePageChange}>
                {#snippet children({pages, currentPage})}
                    <Pagination.Content>
                        <Pagination.Item>
                            <Pagination.Previous/>
                        </Pagination.Item>
                        {#each pages as page (page.key)}
                            {#if page.type === "ellipsis"}
                                <Pagination.Item>
                                    <Pagination.Ellipsis/>
                                </Pagination.Item>
                            {:else}
                                <Pagination.Item>
                                    <Pagination.Link {page} isActive={currentPage === page.value}>
                                        {page.value}
                                    </Pagination.Link>
                                </Pagination.Item>
                            {/if}
                        {/each}
                        <Pagination.Item>
                            <Pagination.Next/>
                        </Pagination.Item>
                    </Pagination.Content>
                {/snippet}
            </Pagination.Root>
        </Tabs.Content>
    {/each}
</Tabs.Root>
