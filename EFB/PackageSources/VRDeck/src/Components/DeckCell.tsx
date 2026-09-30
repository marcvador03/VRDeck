import {TVNode} from "@efb/efb-api";
import {DisplayComponent, FSComponent, ComponentProps, VNode} from "@microsoft/msfs-sdk"
import "./DeckCell.scss";

interface CellProps extends ComponentProps{
    key?: string;
    text: string;
    url: string;
    icontype: string;
};

export class Cell extends DisplayComponent<CellProps> {
    private rootRef = FSComponent.createRef<HTMLDivElement>();
    private imgRef = FSComponent.createRef<HTMLImageElement>(); 
    private labelRef = FSComponent.createRef<HTMLDivElement>();

    public update(text: string, url: string | null, icontype?: string): void {
        this.labelRef.instance.textContent = text;
        this.rootRef.instance.classList.toggle("icon-half", icontype === "half");
        const img = this.imgRef.instance;
        if (url) {
            img.style.display = "";
            img.src = url;
        } else {
            img.style.display = "none";
        }
    }

    public onAfterRender(node: VNode): void {
       this.update(this.props.text, this.props.url || null, this.props.icontype)
    };

    public render(): TVNode<HTMLDivElement> {
        return (
            <div ref = {this.rootRef} class="cell">
                <img ref = {this.imgRef} class="cell-icon" src="" alt="" />
                <div ref = {this.labelRef} class="cell-label">{this.props.text || ""}</div>
            </div>
        )
    };
};