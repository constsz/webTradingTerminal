<script>
    import ActiveTrade from "./activeTrade/ActiveTrade.svelte";

    export let symbolsProps = [];
    export let selectedSymbolCard = 0;

    $: selected = selectedSymbolCard;

    $: ss = symbolsProps;

    let localSs = [];

    $: {
        localSs = [...ss];
        localSs.sort((a, b) => {
            return a.PnlPct - b.PnlPct;
        });
    }
</script>

<main>
    <div class="two-column-cont">
        <div class="side-group">
            <div class="side_header">LONG</div>
            <div class="card-list">
                {#each localSs as symbolProps, i}
                    {#if (symbolProps.PositionStatus == 2 || symbolProps.PositionStatus == 3 || symbolProps.PositionStatus == 0) && symbolProps.Side == 0}
                        <ActiveTrade s={symbolProps} symbolNumber={i} selectedSymbolCard={selected} on:delete-symbol />
                        <br />
                    {/if}
                {/each}
            </div>
        </div>
        <div class="side-group">
            <div class="side_header">SHORT</div>
            <div class="card-list">
                {#each localSs as symbolProps, i}
                    {#if (symbolProps.PositionStatus == 2 || symbolProps.PositionStatus == 3 || symbolProps.PositionStatus == 0) && symbolProps.Side == 1}
                        <ActiveTrade s={symbolProps} symbolNumber={i} selectedSymbolCard={selected} on:delete-symbol />
                        <br />
                    {/if}
                {/each}
            </div>
        </div>
    </div>
</main>

<style>
    .two-column-cont {
        display: flex;
        flex-direction: row;
        margin-left: auto;

        justify-content: right;
    }

    .side_header {
        margin-top: 1rem;
        margin-right: 1rem;

        font-size: 12px;
        font-weight: 800;
        letter-spacing: 2px;

        color: gray;

        display: flex;
        justify-content: right;
    }
    .card-list {
        display: flex;
        margin-left: auto;
        justify-content: right;

        margin-top: 0.5rem;
        margin-right: 1rem;
        margin-bottom: 2rem;

        flex-wrap: wrap;
        width: 560px;
    }
</style>
