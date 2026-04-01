<script lang="ts">
    import * as Pagination from "$lib/components/ui/pagination/index";
    import {goto} from "$app/navigation";
    import HistoryListItem from "$lib/components/top/tracks/HistoryListItem.svelte";
    import ContentListWrapper from "$lib/components/top/ContentListWrapper.svelte";
    import {page} from "$app/state";

    let {data} = $props();

    const handlePageChange = (newPage: number) => {
        const query = new URLSearchParams(page.url.searchParams.toString());
        query.set("page", newPage.toString());
        goto(`?${query.toString()}`, {keepFocus: true});
    };

</script>
<ContentListWrapper title="History">
    <div class="flex flex-col gap-2">
        {#each data.history.results as track}
            <HistoryListItem track={track} parentType="artists"/>
        {/each}
    </div>
</ContentListWrapper>
<Pagination.Root count={data.history.pagination.total_count} perPage={10} page={data.page}
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
