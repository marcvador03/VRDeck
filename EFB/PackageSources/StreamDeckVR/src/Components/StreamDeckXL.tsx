import { GamepadUiView, RequiredProps, TVNode, UiViewProps } from "@efb/efb-api";
import { FSComponent, Subject } from "@microsoft/msfs-sdk";
import "./StreamDeckXL.scss";

interface StreamDeckXLProps extends RequiredProps<UiViewProps, "appViewService" | "bus"> {
  title?: string;
  color?: string;
}
type ConnStatus = "connecting" | "open" | "closed" | "closing";

declare const APP_VERSION: string;

export class StreamDeckXL extends GamepadUiView<HTMLDivElement, StreamDeckXLProps> {
  public readonly tabName = StreamDeckXL.name;
  private cellSubjects: Subject<string>[] = Array.from({ length: 32 }, (_, i) => 
    Subject.create(String(i + 1))
  );

 private status = Subject.create<ConnStatus>("connecting");
 private statusClass = this.status.map((s: ConnStatus) => `sd-status sd-status-${s}`);
 private statusText = this.status.map(s => s === "open" ? "Connected" : s === "connecting" ? "Connecting…" : "Disconnected");

 public onAfterRender(node: TVNode): void {
  console.log("[StreamDeckXL] Setting up bus sub");
  this.props.bus.on("streamdeck-connection-status", (state: ConnStatus) => {
    this.status.set(state);
    console.log("[StreamDeckXL] Status received!", state);
  });
  this.props.bus.on("streamdeck-labels-updated", (json) => {
    console.log("[StreamDeckXL] Event received!", json);
    this.updateLabels(json);
  });
}

  public updateLabels(json: { buttons: { row: number; col: number; label: string }[] }): void {
    console.log("[StreamDeckXL] Update called")
    this.cellSubjects.forEach((sub, i) => sub.set(String("")));
    json.buttons.forEach(button => {
      if (button.row >= 0 && button.row < 4 && button.col >= 0 && button.col < 8) {
        const index = button.row * 8 + button.col;
        this.cellSubjects[index].set(button.label);
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
            <div class="sd-cell" key={`cell-${index}`}>
              {this.cellSubjects[index]}
            </div>
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