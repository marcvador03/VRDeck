import {TVNode} from "@efb/efb-api";
import {DisplayComponent, FSComponent, ComponentProps, VNode} from "@microsoft/msfs-sdk"
import "./DeckCell.scss";

interface CellProps extends ComponentProps{
    key?: string;
    text: string;
    url: string;
    icontype: string;
    keylogic?: string;
};

export class Cell extends DisplayComponent<CellProps> {
    private rootRef = FSComponent.createRef<HTMLDivElement>();
    private imgRef = FSComponent.createRef<HTMLImageElement>(); 
    private labelRef = FSComponent.createRef<HTMLDivElement>();
    private logicRef = FSComponent.createRef<HTMLDivElement>();
    private showingLogic = false;

    private readonly onClick = (): void => {
        this.showingLogic = !this.showingLogic;
        this.rootRef.instance.classList.toggle("show-logic", this.showingLogic);
    };

    public update(text: string, url: string | null, icontype?: string, keylogic?: string): void {
        this.labelRef.instance.textContent = text;
        this.rootRef.instance.classList.toggle("icon-half", icontype === "half");
        const img = this.imgRef.instance;
        if (url) {
            img.style.display = "";
            img.src = url;
        } else {
            img.style.display = "none";
        }
        const logic = this.logicRef.instance;
        logic.innerHTML = "";
        for (const entry of (keylogic ?? "").split("\n")) {
            if (!entry) continue;
            const line = document.createElement("div");
            line.className = "logic-line";
            line.textContent = entry;
            logic.appendChild(line);
        }
        if (!keylogic || keylogic.length === 0) {
            this.rootRef.instance.classList.remove("show-logic");
            this.showingLogic = false;
        }
    };    

    public onAfterRender(node: VNode): void {
       this.update(this.props.text, this.props.url || null, this.props.icontype, this.props.keylogic)
       this.rootRef.instance.addEventListener("click", this.onClick);
    };

    public destroy(): void {
        this.rootRef.instance?.removeEventListener("click", this.onClick);
    }

    public render(): TVNode<HTMLDivElement> {
        return (
            <div ref = {this.rootRef} class="cell">
                <img ref = {this.imgRef} class="cell-icon" src="" alt="" />
                <div ref = {this.labelRef} class="cell-label">{this.props.text || ""}</div>
                <div ref={this.logicRef} class="cell-logic">{this.props.keylogic || ""}</div>
            </div>
        )
    };
};