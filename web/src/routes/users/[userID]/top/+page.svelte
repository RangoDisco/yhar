<script>
    import HistoryList from "$lib/components/top/tracks/HistoryList.svelte";
    import * as Tabs from "$lib/components/ui/tabs/index";
    import TrackList from "$lib/components/top/tracks/TrackList.svelte";
    import ContentGrid from "$lib/components/top/ContentGrid.svelte";

    const {data} = $props();
    const periods = ['week', 'month', 'year', 'overall'];
</script>
<div class="flex flex-col gap-8">
    <Tabs.Root value="week">
        <Tabs.List class="w-full">
            {#each periods as period}
                <Tabs.Trigger value={period}>{period.charAt(0).toUpperCase() + period.slice(1)}</Tabs.Trigger>
            {/each}
        </Tabs.List>
        {#each periods as period}
            <Tabs.Content value={period} class="flex flex-col gap-8">
                {#await data[period]}
                    Loading...
                {:then periodData}
                    <ContentGrid title="Top artists" items={periodData.artists.result} contentType="artists"
                                 url="top/artists"/>
                    <ContentGrid title="Top albums" items={periodData.albums.result} contentType="albums" url="albums"/>
                    <TrackList tracks={periodData.tracks.result}/>
                {/await}
            </Tabs.Content>
        {/each}
    </Tabs.Root>
    <HistoryList tracks={data.history.result} mode="artists"/>
</div>
