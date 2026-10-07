<script lang="ts">
	import ContentCardImage from '#lib/components/top/contentCard/ContentCardImage.svelte';
	import ContentCardWrapper from '#lib/components/top/contentCard/ContentCardWrapper.svelte';
	import { cn } from '#lib/utils.ts';

	type Props = {
		title: string;
		contentID: string | null;
		pictureUrl: string | null;
		contentType: string;
	};

	let { title, contentID = $bindable(null), pictureUrl, contentType }: Props = $props();
</script>

<article class={cn('flex flex-col items-center justify-center gap-2')}>
	<ContentCardWrapper {contentID} {contentType}>
		<ContentCardImage {contentType} {pictureUrl} />
		<p class="mb-auto line-clamp-2 h-12 text-center md:text-lg">{title}</p>
	</ContentCardWrapper>
</article>
