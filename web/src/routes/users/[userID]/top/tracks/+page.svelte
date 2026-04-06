<script lang="ts">
    import * as Tabs from "$lib/components/ui/tabs/index";
    import ContentListItem from "$lib/components/top/ContentListItem.svelte";
    import {handlePeriodChange, PER_PAGE} from "$lib/pagination";
    import Paginator from "$lib/components/Paginator.svelte";

    let {data} = $props();
    const periods = ["week", "month", "year", "overall"];
    let currentPage = $derived(parseInt(data.page));
    let firstIndex = $derived((currentPage - 1) * PER_PAGE + 1);

</script>

<Tabs.Root value={data.period} onValueChange={handlePeriodChange}>
    <Tabs.List class="w-full">
        {#each periods as period}
            <Tabs.Trigger value={period}>{period.charAt(0).toUpperCase() + period.slice(1)}</Tabs.Trigger>
        {/each}
    </Tabs.List>
    {#each periods as period}
        <Tabs.Content value={period} class="flex flex-col gap-8">
            {#each data.tracks.results as track, i (track.id)}
                <ContentListItem index={firstIndex + i} title={track.title} parents={track.artists}
                                 pictureUrl={track.picture_url}
                                 scrobbleCount={track.scrobble_count}
                                 parentType="artists"
                                 contentType="tracks"/>
            {/each}
            <Paginator totalCount={data.tracks.pagination.total_count} page={currentPage}/>
        </Tabs.Content>
    {/each}
</Tabs.Root>
