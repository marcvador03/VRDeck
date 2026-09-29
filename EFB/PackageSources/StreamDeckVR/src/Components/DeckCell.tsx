import {TVNode} from "@efb/efb-api";
import {DisplayComponent, FSComponent, ComponentProps, VNode} from "@microsoft/msfs-sdk"
import "./DeckCell.scss";

interface CellProps extends ComponentProps{
    key?: string;
    text: string;
    url: string;
};

export class Cell extends DisplayComponent<CellProps> {
    private imgRef = FSComponent.createRef<HTMLImageElement>(); 
    private labelRef = FSComponent.createRef<HTMLDivElement>();

    public update(text: string, url: string | null): void {
        this.labelRef.instance.textContent = text;
        const img = this.imgRef.instance;
        if (url) {
            img.style.display = "";
            img.src = url;
        } else {
            img.style.display = "none";
        }
    }

    public onAfterRender(node: VNode): void {
       this.update(this.props.text, this.props.url || null)
    };

    public render(): TVNode<HTMLDivElement> {
        return (
            <div class="cell">
                <img ref = {this.imgRef} class="cell-icon" src="" alt="" />
                <div ref = {this.labelRef} class="cell-label">{this.props.text || ""}</div>
            </div>
        )
    };
};