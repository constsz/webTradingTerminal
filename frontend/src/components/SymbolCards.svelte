<script>
    import SymbolCard from "./symbolCard/SymbolCard.svelte";

    export let symbolsProps = [];
    export let selectedSymbolCard = 0;

    $: selected = selectedSymbolCard;

    $: ss = symbolsProps;
</script>

<main>
    <div class="card-list">
        {#each ss as symbolProps, i}
            <SymbolCard s={symbolProps} symbolNumber={i} selectedSymbolCard={selected} on:delete-symbol />
            <br />
        {/each}
    </div>
</main>

<style>
    .card-list {
        display: flex;
        justify-content: right;
        margin-right: 1rem;
        margin-bottom: 2rem;

        flex-wrap: wrap;
    }
</style>
