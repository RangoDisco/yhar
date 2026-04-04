<script lang="ts">
    import {page} from "$app/state";
    import {Button} from "$lib/components/ui/button";
</script>

{#if page.error}
    <div class="container mx-auto flex items-center justify-center">
        {#if page.status === 404}
            <div class="w-md flex flex-col gap-8 ">
                <p class="text-xl">Oops, the content you're looking for doesn't exist</p>
                <Button href="/">Back to my stats</Button>
            </div>
        {:else}
            Oops, something went wrong. Please try again later
        {/if}
    </div>
{/if}
