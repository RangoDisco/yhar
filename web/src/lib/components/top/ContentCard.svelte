<script lang="ts">
    import * as Avatar from "$lib/components/ui/avatar/index";

    let {
        title,
        pictureUrl,
        contentType
    } = $props();

</script>

<article class="flex flex-col items-center justify-center gap-2 w-[30%] md:w-[15%]">
    <img class="{contentType !== 'artists' ? 'rounded-md' : 'rounded-full'} w-full aspect-square" src={pictureUrl}
         alt={`${title}'s picture`}/>
    <p class="text-center mb-auto">{title}</p>
</article>
