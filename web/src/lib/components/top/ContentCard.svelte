<script lang="ts">
    import * as Avatar from "$lib/components/ui/avatar/index";

    let {
        title,
        pictureUrl,
        contentType
    } = $props();

</script>

<article class="flex flex-col items-center justify-center">
    <Avatar.Root class="{contentType !== 'artists' ? 'rounded-md' : ''} h-[80%] w-[80%] aspect-square">
        <Avatar.Image src={pictureUrl}
                      alt={`${title}'s picture`}/>
        <Avatar.Fallback>{title}</Avatar.Fallback>
    </Avatar.Root>
    <p>{title}</p>
</article>
