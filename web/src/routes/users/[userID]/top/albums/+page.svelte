<script lang="ts">
    import * as Tabs from "$lib/components/ui/tabs/index";
    import * as Pagination from "$lib/components/ui/pagination/index";
    import ContentListItem from "$lib/components/top/ContentListItem.svelte";
    import {goto} from "$app/navigation";
    import {page} from "$app/state";

    let {data} = $props();
    const periods = ["week", "month", "year", "overall"];

    const handlePeriodChange = (period: string) => {
        page.url.searchParams.set("period", period);
        goto(page.url, {keepFocus: true});
    };

    const handlePageChange = (newPage: number) => {
        page.url.searchParams.set("page", newPage.toString());
        goto(page.url);
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
            {#each data.albums.result as album, i}
                <ContentListItem index={i} contentID={album.id} title={album.title}
                                 pictureUrl={album.picture_url}
                                 scrobbleCount={album.scrobble_count}
                                 parentType="artists"
                                 parents={album.artists}
                                 contentType="albums"/>
            {/each}
            <Pagination.Root count={data.albums.pagination.total_count} perPage={10} page={data.page}
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
