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
            {#each data.artists.results as artist, i (artist.id)}
                <ContentListItem index={firstIndex + i} contentID={artist.id} title={artist.name} parents={[]}
                                 pictureUrl={artist.picture_url}
                                 scrobbleCount={artist.scrobble_count}
                                 parentType={null}
                                 contentType="artists"/>
            {/each}
            <Paginator totalCount={data.artists.pagination.total_count} page={currentPage}/>
        </Tabs.Content>
    {/each}
</Tabs.Root>
