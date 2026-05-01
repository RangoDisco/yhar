<script lang="ts">
    import * as Pagination from "$lib/components/ui/pagination/index";
    import Paginator from "$lib/components/Paginator.svelte";
    import HistoryList from "$lib/components/top/tracks/HistoryList.svelte";

    let {data} = $props();

    let currentPage = $derived(parseInt(data.page));

</script>
<HistoryList scrobbles={data.history.results} url={null} parentType="artists"/>
<Paginator totalCount={data.history.pagination.total_count} page={currentPage}/>
