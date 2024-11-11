<script>
    import { createEventDispatcher } from "svelte";
    import { baseUrl } from "../baseUrl";

    const dispatch = createEventDispatcher();

    export let s = {};
    export let symbolNumber = 0;
    export let selectedSymbolCard = 0;

    $: selected = selectedSymbolCard;

    // Dynamic variables: Side
    let sideString;
    let sideColor;
    let sideTextColor;

    const sideStringEnum = {
        0: "LONG",
        1: "SHORT",
    };

    const entryTypeEnum = {
        0: "EntryByPrice",
        1: "EntryByPercent",
    };

    const setSide = (sideInt) => {
        console.log(`symbolNumber: ${symbolNumber} | selected: ${selected}`);

        s.Side = sideInt;
        dynamicSide();
        console.log("setSide: s.Side == " + s.Side);
    };

    const setEntryType = (entryTypeInt) => {
        s.EntryType = entryTypeInt;
    };

    const dynamicSide = () => {
        switch (s.Side) {
            case 0:
                sideString = sideStringEnum[0];
                sideColor = "#1b733d";
                sideTextColor = "#6bd193";
                break;
            case 1:
                sideString = sideStringEnum[1];
                sideColor = "#791b51";
                sideTextColor = "#d16ba9";
                break;
        }
    };

    const positionStatusColorDynamic = (positionStatus, pnl) => {
        if (positionStatus == 0) {
            return "#485366";
        } else if (positionStatus == 1) {
            return "#746C00";
        } else {
            if (pnl > 0) {
                return "#1b733d";
            } else {
                return "#791b51";
            }
        }
    };

    $: positionStatusColor = positionStatusColorDynamic(s.PositionStatus, s.PnlPct);

    // Dynamic variables: Mode
    let modeString;
    $: modeColor = "#4444";

    const dynamicMode = () => {
        switch (s.Mode) {
            case 0:
                modeString = "MANUAL";
                modeColor = "#14347d";
                break;
            case 1:
                modeString = "AUTO";
                modeColor = "#3f1f73";
                break;
        }
    };

    // Dynamic variables: PositionStatus

    const positionStatusEnum = {
        0: "StatusNotInPositionNoOpenOrders",
        1: "StatusNotInPositionOrdersPlaced",
        2: "StatusInPositionNoOpenOrders",
        3: "StatusInPositionAddOrderPlaced",
    };

    const plusOrMinus = (n) => {
        return n < 0 ? " " : " +";
    };
    const _ = "\xa0\xa0\xa0\xa0";

    const dynamicPositionStatus = (s) => {
        switch (s.PositionStatus) {
            case 0:
                return "Blank";
            case 1:
                return `${sideStringEnum[s.Side]} Order Placed`;
            case 2:
                return `${sideStringEnum[s.Side]}: ${s.NumOfEntries}($${s.BaseQuantityUsdt.toFixed(0)}) ${_} ${s.AvgPrice} ${_} ${plusOrMinus(
                    s.PnlPct
                )}$${s.PnlUsd.toFixed(2)} ${_} ${plusOrMinus(s.PnlPct)}${s.PnlPct.toFixed(2)}%`;
            case 3:
                return `${sideStringEnum[s.Side]}: ${s.NumOfEntries}($${s.BaseQuantityUsdt.toFixed(0)}) ${_} ${s.AvgPrice} ${_} ${plusOrMinus(
                    s.PnlPct
                )}$${s.PnlUsd.toFixed(2)} ${_} ${plusOrMinus(s.PnlPct)}${s.PnlPct.toFixed(2)}% ${_} | ${_} +Order`;
        }
    };

    $: positionStatusString = dynamicPositionStatus(s);

    dynamicSide();
    dynamicMode();

    // ----------------------------------------

    const deleteSymbol = async (symbolName) => {
        let msgBody = {
            Symbol: symbolName,
        };

        const url = `${baseUrl}/symbolSettings/delete`;

        await fetch(url, {
            headers: {
                Accept: "application/json",
                "Content-Type": "application/json",
            },
            method: "POST",
            body: JSON.stringify(msgBody),
        });

        dispatch("delete-symbol", symbolName);
    };

    let entrySideText = {
        0: "$",
        1: "%",
    };
    const toggleEntryType = () => {
        s.EntryType = s.EntryType == 0 ? 1 : 0;

        switch (s.EntryType) {
            case 0:
                s.Entry = 0;
                break;
            case 1:
                s.Entry = 0.1;
                break;

            default:
                break;
        }
    };

    const toggleSide = () => {
        if (s.PositionStatus != 2 && s.PositionStatus != 3) {
            console.log("s.Side: " + s.Side);
            s.Side = s.Side == 0 ? 1 : 0;
            dynamicSide();
            console.log("s.Side: " + s.Side);
        }
    };

    // ----------------------------------------

    const ActionType = {
        SetOrder: 0,
        CancelAllOrders: 1,
        MarketEntry: 2,
        MarketExit: 3,
    };

    // prepareSymbolSettingsUpdate is generating msgBody for API call
    // receiving paramenter: actionType
    // actionType is int coming from ActionType.SetOrder etc
    const symbolActionCall = async (actionType) => {
        s.Entry = Number(s.Entry);
        s.TP = Number(s.TP);
        s.SL = Number(s.SL);

        let msgBody = {
            ActionType: actionType,
            SymbolSettings: s,
        };

        const url = `${baseUrl}/symbolAction`;

        let validated = false;

        if (s.SL > 0 && s.TP > 0) {
            validated = true;
        }

        if (validated) {
            await fetch(url, {
                headers: {
                    Accept: "application/json",
                    "Content-Type": "application/json",
                },
                method: "POST",
                body: JSON.stringify(msgBody),
            });
        }
    };

    const actionSetOrder = async () => {
        await symbolActionCall(ActionType.SetOrder);
    };

    const actionCancelAllOrders = async () => {
        await symbolActionCall(ActionType.CancelAllOrders);
    };

    const actionMarketEntry = async () => {
        await symbolActionCall(ActionType.MarketEntry);
    };

    const actionMarketExit = async () => {
        await symbolActionCall(ActionType.MarketExit);
    };
</script>

<main>
    <div class="card-container">
        <div class="fields">
            <div class="field-A">
                <div class="menu-row">
                    <div class="symbol-name" style="--mode-color: {modeColor}">
                        <div class="text-center">{s.Symbol}</div>
                    </div>
                    <div class="side" style="--side-color: {sideColor}; --side-text-color: {sideTextColor}">
                        <div class="text-center" on:click={toggleSide}>
                            {sideString}
                        </div>
                    </div>
                </div>

                <div class="field-entry">
                    <div class="field-entry-col field-entry-text">
                        <div class="field-entry-row">Entry</div>
                        <div class="field-entry-row">SL TP</div>
                    </div>

                    <div class="field-entry-col field-entry-col-input">
                        <div class="field-entry-row">
                            <form on:submit|preventDefault={actionSetOrder}>
                                <input class="input" type="text" bind:value={s.Entry} id={"input_Entry_" + symbolNumber} />
                                <div class="entry-type-toggle" on:click={toggleEntryType}>{entrySideText[s.EntryType]}</div>
                            </form>
                        </div>
                        <div class="field-entry-row">
                            <form on:submit|preventDefault={actionSetOrder}>
                                <input class="input" type="text" style="max-width: 34px" bind:value={s.SL} id={"input_SL_" + symbolNumber} />
                            </form>
                            <form on:submit|preventDefault={actionSetOrder}>
                                <input class="input" type="text" style="max-width: 34px; margin-left: 2px" bind:value={s.TP} id={"input_TP_" + symbolNumber} />
                            </form>
                        </div>
                    </div>
                </div>
            </div>

            <div class="field-B">
                <div class="menu-row gap-1px">
                    <div class="mode" style="--mode-color: {modeColor}">
                        <div class="text-center" style="opacity: 0.6">
                            <div class="mode-inner-text" />
                        </div>
                    </div>
                    <div
                        on:click={() => {
                            deleteSymbol(s.Symbol);
                        }}
                        class="close-btn"
                    >
                        <div class="text-center" style="opacity: 0.4">x</div>
                    </div>
                </div>
                <div class="row-manual gap-1px">
                    <button class="btn-manual cursor-add" on:click={actionSetOrder} id={"btn_SetOrder_" + symbolNumber}>Set Order</button>
                    <button class="btn-manual btn-manual-small cursor-delete" on:click={actionCancelAllOrders} id={"btn_CancelAllOrders_" + symbolNumber}
                        >x</button
                    >
                </div>
                <div class="row-manual gap-1px">
                    <button class="btn-manual cursor-add" on:click={actionMarketEntry} id={"btn_MarketEntry_" + symbolNumber}>Market Entry</button>
                    <button class="btn-manual btn-manual-small cursor-delete" on:click={actionMarketExit} id={"btn_MarketExit_" + symbolNumber}>x</button>
                </div>
            </div>
        </div>
        <div class="position-status" style="--position-status-color: {positionStatusColor}"><div class="position-status-text">{positionStatusString}</div></div>
    </div>
</main>

<style>
    .card-container {
        font-size: 14px;

        display: flex;
        flex-direction: column;
        width: 330px;
        height: 126px;
        background: #181b24;
        margin: 1rem 0rem 0rem 1rem;
    }

    /* FIELDS */
    .fields {
        display: flex;
        width: 100%;
        height: 100%;
    }

    .field-A {
        display: flex;
        flex-direction: column;
        justify-content: space-between;

        flex-basis: 70%;
        margin-right: 1px;

        height: auto;
    }

    .menu-row {
        display: flex;
    }

    .gap-1px {
        gap: 1px;
    }

    .symbol-name {
        background: var(--mode-color);

        flex-basis: 55%;

        font-weight: bold;
        letter-spacing: 1.6px;

        padding-left: 12px;
        height: 24px;
    }

    .side {
        background: var(--side-color);
        flex-basis: 45%;

        text-align: center;

        font-size: 11px;
        font-weight: bold;
        letter-spacing: 2.4px;
        color: var(--side-text-color);

        margin-right: 1px;
        height: 24px;

        user-select: none;
        cursor: pointer;
    }

    .field-entry {
        display: flex;
        margin-top: 8px;
    }

    .field-entry-col {
        display: flex;
        flex-direction: column;
        flex-basis: 30%;
    }

    .field-entry-col-input {
        flex-basis: 70%;
    }

    .field-entry-row {
        height: 34px;
    }

    .field-entry-text {
        padding-top: 6px;
        padding-left: 12px;
        user-select: none;
    }

    form {
        display: inline;
    }

    .position-status {
        font-size: 11px;
        font-weight: bold;
        padding-left: 12px;
        color: #fff;
        background: var(--position-status-color);
        height: 20px;
        width: auto;
        user-select: none;
    }

    .position-status-text {
        transform: translateY(3px);
        letter-spacing: 0.75px;
    }

    .entry-type-toggle {
        color: #858585;
        font-size: 12px;
        font-weight: bold;
        /* display: inline-block; */
        display: flex;
        position: relative;
        text-align: center;
        align-items: center;
        justify-content: center;

        left: 80px;
        top: -28px;
        background-color: rgb(42, 45, 56);
        padding: 5px;
        border-radius: 3px;
        cursor: pointer;
        height: 16px;
        width: 18px;
    }

    .input {
        max-width: 100px;
        height: 28px;
        background: #131722;
        padding-left: 8px;
        border: 1px solid #464c5e;
        border-radius: 4px;
        color: #fefefe;
    }

    .input::placeholder {
        opacity: 1;
        color: #fefefe;
    }

    /*  FIELD B */
    .field-B {
        flex-basis: 46.9%;
        display: flex;
        flex-direction: column;
        gap: 1px;
        /* flex-wrap: wrap; */
        justify-content: space-between;
        align-content: space-between;
    }

    .mode {
        flex-basis: 72%;

        background: var(--mode-color);
        font-size: 11px;
        font-weight: 500;
        letter-spacing: 2.4px;
    }

    .mode-inner-text {
        padding-left: 12px;
    }

    .close-btn {
        flex-basis: 28%;

        background: #2b2a35;
        text-align: center;

        height: 24px;

        cursor: pointer;

        border: none;
    }

    .row-manual {
        display: flex;
        height: 100%;
    }

    .btn-manual {
        flex-basis: 72%;
        background-color: #4444;
        border: none;

        font-weight: bold;
        color: #919191;

        cursor: crosshair;
    }

    .cursor-add {
        cursor: cell;
    }

    .cursor-delete {
        cursor: not-allowed;
    }

    .btn-manual-small {
        flex-basis: 28%;
    }

    .text-center {
        position: relative;
        top: 50%;
        transform: translateY(-50%);
    }
</style>
