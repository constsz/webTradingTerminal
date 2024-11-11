<script>
    import { onMount } from "svelte";
    import SymbolCards from "./components/SymbolCards.svelte";
    import ActiveTrades from "./components/ActiveTrades.svelte";
    import Menu from "./components/Menu.svelte";
    import { baseUrl } from "./components/baseUrl";
    import ActiveTrade from "./components/activeTrade/ActiveTrade.svelte";

    // ----------------------------------------
    // KEYBOARD SHORTCUTS

    $: selectedSymbolCard = 0;
    let textMode = true;

    const isNaturalKey = (key) => {
        let listOfNaturalKeys = ["backspace", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "."];

        for (let i in listOfNaturalKeys) {
            if (key == listOfNaturalKeys[i]) {
                return true;
            }
        }

        return false;
    };

    document.addEventListener("keydown", (e) => {
        let id = "";

        const key = e.key.toLowerCase();

        switch (key) {
            // Text mode
            case "`":
                if (pageSelect.cards == true) {
                    pageSelect.cards = false;
                    pageSelect.activeTrades = true;
                } else {
                    pageSelect.cards = true;
                    pageSelect.activeTrades = false;
                }
        }
    });

    // ----------------------------------------
    // INITIAL LOADING of SymbolSettings

    let symbolsProps = [];

    onMount(async () => {
        try {
            const res = await fetch(`${baseUrl}/symbolSettings/get`);
            symbolsProps = await res.json();
        } catch (e) {
            error = e;
        }
    });

    // ----------------------------------------
    // POLLING of Symbol Status

    const setupPoller = () => {
        // symbolsProps = setInterval(doPoll(), 1000);
        setInterval(() => {
            doPoll();
        }, 1000);
    };

    const doPoll = async () => {
        const res = await fetch(`${baseUrl}/symbolsStatus`);
        let symbolsStatus = await res.json();

        for (const i in symbolsProps) {
            let symbol = symbolsProps[i].Symbol;

            symbolsProps[i].PositionStatus = symbolsStatus[symbol].PositionStatus;
            symbolsProps[i].BaseQuantityUsdt = symbolsStatus[symbol].BaseQuantityUsdt;
            symbolsProps[i].NumOfEntries = symbolsStatus[symbol].NumOfEntries;
            symbolsProps[i].AvgPrice = symbolsStatus[symbol].AvgPrice;
            symbolsProps[i].PnlUsd = symbolsStatus[symbol].PnlUsd;
            symbolsProps[i].PnlPct = symbolsStatus[symbol].PnlPct;

            if (symbolsStatus[symbol].PositionStatus > 1) {
                symbolsProps[i].Side = symbolsStatus[symbol].Side;
            }
        }
    };

    $: setupPoller();

    // ----------------------------------------
    // FUNCTIONS

    let symbolExists = false;

    const symbolExistsWarning = async () => {
        symbolExists = true;
        setTimeout(() => {
            symbolExists = false;
        }, 3000);
    };

    const addSymbol = async (e) => {
        const symbol = e.detail;

        for (const i in symbolsProps) {
            if (symbol == symbolsProps[i].Symbol) {
                symbolExistsWarning();
            }
        }

        if (!symbolExists) {
            const msgBody = {
                Symbol: symbol,
            };

            const url = `${baseUrl}/symbolSettings/add`;

            await fetch(url, {
                headers: {
                    Accept: "application/json",
                    "Content-Type": "application/json",
                },
                method: "POST",
                body: JSON.stringify(msgBody),
            });

            const newSymbolCard = {
                Symbol: symbol,
                Mode: 0,
                Side: 0,

                Entry: 0.1,
                EntryType: 1,
                SL: 0.2,
                TP: 0.25,

                PositionStatus: 0,
                BaseQuantityUsdt: 0,
                NumOfEntries: 0,
                AvgPrice: "",
                PnlUsd: 0.0,
                PnlPct: 0.0,

                TF: 0,
                ATRM: 0,
                Execution: 0,
                BotStatus: 0,
            };

            symbolsProps = [...symbolsProps, newSymbolCard];

            setTimeout(() => {
                window.scrollBy(0, document.body.scrollHeight);
            }, 5);
        }
    };

    const deleteSymbol = async (e) => {
        const symbol = e.detail;

        symbolsProps = symbolsProps.filter((value, index, arr) => {
            return value.Symbol != symbol;
        });
    };

    let pageSelect = {
        cards: true,
        activeTrades: false,
    };

    const menuSelectPage = (e) => {
        let page = e.detail;
        switch (page) {
            case "cards":
                pageSelect.cards = true;
                pageSelect.activeTrades = false;
                break;
            case "activeTrades":
                pageSelect.cards = false;
                pageSelect.activeTrades = true;
                break;
            default:
                break;
        }
    };
</script>

<main>
    <div class="board">
        {#if pageSelect.cards}
            <SymbolCards {symbolsProps} {selectedSymbolCard} on:delete-symbol={deleteSymbol} />

            {#if symbolExists}
                <div class="verification-hint">Symbol already exists.</div>
            {/if}
        {/if}
        {#if pageSelect.activeTrades}
            <ActiveTrades {symbolsProps} {selectedSymbolCard} on:delete-symbol={deleteSymbol} />
        {/if}
        <Menu on:add-symbol={addSymbol} {pageSelect} on:menu-selected-page={menuSelectPage} />
    </div>
</main>

<style>
    .board {
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        min-height: 100vh;

        /* background-color: #9ceb95; */
    }

    .verification-hint {
        position: absolute;
        transform: translate(1rem, 50vh);
        color: white;
        background-color: red;
        padding: 28px;
        border-radius: 8px;
        opacity: 0.9;
        box-shadow: 20px;
    }
</style>
