<script lang="ts">
    import disc from "$lib/assets/placeholders/disc-3.svg";
    import user from "$lib/assets/placeholders/user.jpg";

    let {
        contentType,
        pictureUrl,
    } = $props();

    let fallback = $derived(contentType === "artists" ? user : disc);
</script>

<img
    class="w-full aspect-square object-cover {contentType === 'artists' ? 'rounded-full' : 'rounded-md'}"
    src={pictureUrl ?? fallback}
    alt="Content's cover"
    onerror={() => pictureUrl = fallback}
/>
