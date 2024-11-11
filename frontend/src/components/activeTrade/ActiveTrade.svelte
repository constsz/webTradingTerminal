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

    const dynamicSide = () => {
        switch (s.Side) {
            case 0:
                sideString = sideStringEnum[0];
                sideColor = "#0e0e0e";
                sideTextColor = "#6bd193";
                break;
            case 1:
                sideString = sideStringEnum[1];
                sideColor = "#0e0e0e";
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
            if (pnl > 2) {
                return "#0F7336";
            } else if (pnl > 1) {
                return "#105B2D";
            } else if (pnl > 0) {
                return "#0F4725";
            } else if (pnl < -2) {
                return "#830F50";
            } else if (pnl < -1) {
                return "#681142";
            } else if (pnl < 0) {
                return "#58133A";
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
        return n < 0 ? " -" : " +";
    };
    const _ = "\xa0\xa0\xa0\xa0";

    const dynamicPositionStatus = (s) => {
        switch (s.PositionStatus) {
            case 0:
                return `${_} Blank ${_} ${plusOrMinus(s.PnlPct)}${Math.abs(s.PnlPct).toFixed(2)}%`;
            case 1:
                return `${sideStringEnum[s.Side]} Order Placed`;
            case 2:
                return `${_} ${plusOrMinus(s.PnlUsd)}$${Math.abs(s.PnlUsd).toFixed(2)} ${_} ${s.NumOfEntries}($${s.BaseQuantityUsdt.toFixed(0)}) ${_} ${
                    s.AvgPrice
                } ${_}  `;
            case 3:
                return `${_} ${plusOrMinus(s.PnlUsd)}$${Math.abs(s.PnlUsd).toFixed(2)} ${_} ${s.AvgPrice} ${_} ${plusOrMinus(s.PnlPct)}$${s.PnlUsd.toFixed(
                    2
                )} ${_} | +Order`;
        }
    };

    const dynamicPnLString = (s) => {
        switch (s.PositionStatus) {
            case 0:
                return "";
            case 1:
                return "";
            case 2:
                return `${plusOrMinus(s.PnlPct)}${Math.abs(s.PnlPct).toFixed(2)}%`;
            case 3:
                return `${plusOrMinus(s.PnlPct)}${Math.abs(s.PnlPct).toFixed(2)}%`;
        }
    };

    $: positionStatusString = dynamicPositionStatus(s);
    $: positionPnLString = dynamicPnLString(s);

    dynamicSide();
    dynamicMode();

    // ----------------------------------------

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
            <div class="symbol-name" style="--position-status-color: {positionStatusColor}">
                <div class="text-center symbol-inner">{s.Symbol}</div>
            </div>

            <div class="position-status flex-grow" style="--position-status-color: {positionStatusColor}">
                <div class="position-status-text width-80">
                    {positionPnLString}
                </div>

                <div class="side" style="--side-color: {sideColor}; --side-text-color: {sideTextColor}">
                    <div class="text-center">
                        {sideString}
                    </div>
                </div>

                <div class="position-status-text">
                    <div class="position-status-minitext">{positionStatusString}</div>
                </div>
            </div>
        </div>
    </div>
</main>

<style>
    .card-container {
        font-size: 14px;

        display: flex;
        width: 540px;
        height: 24px;
        background: #0e0e0e;
        margin: 0.4rem 0rem 0rem 1rem;
    }

    /* FIELDS */
    .fields {
        display: flex;
        width: 100%;
        height: 100%;
    }

    .position-status {
        display: flex;
        flex-wrap: wrap;
        align-items: center;
        font-weight: 600;
        padding-left: 12px;
        color: #fff;
        background: var(--position-status-color);

        user-select: none;
    }

    .width-80 {
        width: 80px;
    }

    .flex-grow {
        flex-grow: 1;
    }

    .position-status-minitext {
        display: inline;
        font-size: 12px;
        font-weight: 400;

        opacity: 0.6;
    }
    .position-status-text {
        letter-spacing: 0.75px;
    }

    .symbol-name {
        /* background: var(--mode-color); */
        /* background: var(--position-status-color); */
        background-color: var(--position-status-color);
        width: 90px;

        font-weight: bold;
        letter-spacing: 1.6px;

        user-select: none;
    }
    .symbol-inner {
        background-color: rgba(0, 0, 0, 0.2);
        display: flex;
        height: 100%;
        padding-left: 12px;
        align-items: center;
    }

    .side {
        background: rgba(0, 0, 0, 0.6);
        width: 60px;
        height: 18px;

        text-align: center;

        border-radius: 4px;

        font-size: 9px;
        font-weight: bold;
        letter-spacing: 2.4px;
        /* color: var(--side-text-color); */
        color: rgba(255, 255, 255, 0.7);

        margin-right: 1px;

        user-select: none;
    }

    .text-center {
        position: relative;
        top: 50%;
        transform: translateY(-50%);
    }
</style>
