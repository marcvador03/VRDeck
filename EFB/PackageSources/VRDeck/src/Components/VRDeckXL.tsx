import { GamepadUiView, RequiredProps, TVNode, UiViewProps } from "@efb/efb-api";
import { FSComponent, Subject } from "@microsoft/msfs-sdk";
import { Cell } from "./DeckCell";
import "./VRDeckXL.scss";

interface VRDeckXLProps extends RequiredProps<UiViewProps, "appViewService" | "bus"> {
  title?: string;
  color?: string;
}

type ConnStatus = "connecting" | "open" | "closed" | "closing";

declare const APP_VERSION: string;

export class VRDeckXL extends GamepadUiView<HTMLDivElement, VRDeckXLProps> {
  public readonly tabName = VRDeckXL.name;
  private cellSubjects: Subject<string>[] = Array.from({ length: 32 }, (_, i) => 
    Subject.create(String(i + 1))
  );

 private status = Subject.create<ConnStatus>("connecting");
 private statusClass = this.status.map((s: ConnStatus) => `sd-status sd-status-${s}`);
 private statusText = this.status.map(s => s === "open" ? "Connected" : s === "connecting" ? "Connecting…" : "Disconnected");

 private cellContent = Array.from( {length: 32}, () => FSComponent.createRef<Cell>());

 public onAfterRender(node: TVNode): void {
  console.log("[VRDeck] Setting up bus sub");
  this.props.bus.on("streamdeck-connection-status", (state: ConnStatus) => {
    this.status.set(state);
    console.log("[VRDeck] Status received!", state);
  });
  this.props.bus.on("streamdeck-labels-updated", (json) => {
    console.log("[VRDeck] Event received!", json);
    this.updateLabels(json);
  });
}

 public updateLabels(json: { buttons: { row: number; col: number; label: string; icon?: string, icontype?: string, keylogic?: string}[] }): void {
    console.log("[VRDeck] Update called")
    this.cellContent.forEach((ref) => ref.instance.update("", null, undefined));
    json.buttons.forEach(button => {
      if (button.row >= 0 && button.row < 4 && button.col >= 0 && button.col < 8) {
        const index = button.row * 8 + button.col;
        this.cellContent[index].instance.update(button.label, button.icon ?? null, button.icontype, button.keylogic);
      }
    });
  }

 public render(): TVNode<HTMLDivElement> {
    return (
      <div ref={this.gamepadUiViewRef} class="sd-parent">
        <div class="sd-header">
          VRDeck
        </div>
        <div class="sd-container">
          {Array.from({ length: 32 }).map((_, index) => (
            <Cell ref={this.cellContent[index]} key={`cell-${index}`} text="" url="" icontype=""/>
        ))}
        </div>
        <div class="sd-bottom">
          <span class={this.statusClass}></span>
          <span>{this.statusText}</span>
          <span>vv{APP_VERSION}</span>
        </div>
      </div>
    );
  }
}