<script lang="ts">
    import HistoryListItem from "$lib/components/top/tracks/HistoryListItem.svelte";
</script>

<section class="flex flex-col gap-4">
    <div class="flex justify-between items-center">
        <h1 class="text-2xl">History</h1>
        <a class="text-sm text-muted-foreground" href="history">View all</a>
    </div>
    <div class="flex flex-col gap-4">
        <HistoryListItem title="wake up" album="sace6" scrobbledAt="4 minutes ago"
                         pictureUrl="https://wrapped.nero.rangodisco.eu/cacheimages/4780a4842ce999779dee6863c1e3a0fe"
        />
        <HistoryListItem title="Venturing" album="sace6" scrobbledAt="4 minutes ago"
                         pictureUrl="https://wrapped.nero.rangodisco.eu/cacheimages/ea86e05166610135d84e5724be5a23e1"
        />
        <HistoryListItem title="Ptite soeur" album="sace6" scrobbledAt="4 minutes ago"
                         pictureUrl="https://wrapped.nero.rangodisco.eu/cacheimages/4780a4842ce999779dee6863c1e3a0fe"
        />
        <HistoryListItem title="Ptite soeur" album="sace6" scrobbledAt="4 minutes ago"
                         pictureUrl="https://wrapped.nero.rangodisco.eu/cacheimages/4780a4842ce999779dee6863c1e3a0fe"
        />
        <HistoryListItem title="Ptite soeur" album="sace6" scrobbledAt="4 minutes ago"
                         pictureUrl="https://wrapped.nero.rangodisco.eu/cacheimages/4780a4842ce999779dee6863c1e3a0fe"
        />
        <HistoryListItem title="Ptite soeur" album="sace6" scrobbledAt="4 minutes ago"
                         pictureUrl="https://wrapped.nero.rangodisco.eu/cacheimages/4780a4842ce999779dee6863c1e3a0fe"
        />
    </div>
</section>
