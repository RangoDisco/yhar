<script lang="ts">
    import * as Pagination from "$lib/components/ui/pagination/index";
    import HistoryListItem from "$lib/components/top/tracks/HistoryListItem.svelte";
    import ContentListWrapper from "$lib/components/top/ContentListWrapper.svelte";
    import Paginator from "$lib/components/Paginator.svelte";

    let {data} = $props();

    let currentPage = $derived(parseInt(data.page));

</script>
<ContentListWrapper title="History">
    <div class="flex flex-col gap-2">
        {#each data.history.results as track}
            <HistoryListItem track={track} parentType="artists"/>
        {/each}
    </div>
</ContentListWrapper>
<Paginator totalCount={data.history.pagination.total_count} page={currentPage}/>
