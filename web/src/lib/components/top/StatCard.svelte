<script lang="ts">
    import * as Card from "#lib/components/ui/card/index.js";

    type Props = {
        title: string;
        description: string;
    }

    const {title, description}: Props = $props();

</script>

<Card.Root>
    <Card.Content>
        <!--Dump swap bc it matches the style better IMO-->
        <Card.Description class="h-10">{title}</Card.Description>
        <Card.Title class="text-xl">{description}</Card.Title>
    </Card.Content>
</Card.Root>

