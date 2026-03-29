<script lang="ts">
    import disc from "$lib/assets/placeholders/disc-3.svg";
    import user from "$lib/assets/placeholders/user.svg"

    let {
        contentType,
        pictureUrl,
    } = $props();
</script>

{#if contentType === "artists"}
    <img class="rounded-full w-full aspect-square" src={pictureUrl ?? user}
         alt={`Content picture`}/>
{:else}
    <img class="rounded-md w-full aspect-square" src={pictureUrl ?? disc}
         alt={`Content picture`}/>
{/if}
