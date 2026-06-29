import { GamepadUiView, RequiredProps, TVNode, UiViewProps } from "@efb/efb-api";
import { FSComponent } from "@microsoft/msfs-sdk";
import "./StreamDeckXL.scss";

interface StreamDeckXLProps extends RequiredProps<UiViewProps, "appViewService" | "bus"> {
  title?: string;
  color?: string;
}

export class StreamDeckXL extends GamepadUiView<HTMLDivElement, StreamDeckXLProps> {
  public readonly tabName = StreamDeckXL.name;
  private cellLabels: string[][] = Array.from({ length: 4 }, (_, row) =>
    Array.from({ length: 8 }, (_, col) => String(row * 8 + col + 1))
  );

 public onAfterRender(node: TVNode): void {
   console.log("[StreamDeckXL] Setting up bus sub");
    this.props.bus.on("streamdeck-labels-updated", (json) => {
      console.log("[StreamDeckXL] Event received!", json);
      this.updateLabels(json);
    });
  }

  public updateLabels(json: { buttons: { row: number; col: number; label: string }[] }): void {
    console.log("[StreamDeckXL] Update called")
    this.cellLabels = Array.from({ length: 4 }, (_, row) =>
      Array.from({ length: 8 }, (_, col) => String(row * 8 + col + 1))
    );
    json.buttons.forEach(button => {
      if (button.row >= 0 && button.row < 4 && button.col >= 0 && button.col < 8) {
        this.cellLabels[button.row][button.col] = button.label;
      }
    });
    this.props.appViewService.update
  }

  public render(): TVNode<HTMLDivElement> {
    return (
      <div ref={this.gamepadUiViewRef} class="streamdeck-container">
        {Array.from({ length: 4 }).map((_, row) =>
          Array.from({ length: 8 }).map((_, col) => (
            <div key={`cell-${row}-${col}`} class="streamdeck-cell">
              {this.cellLabels[row][col]}
            </div>
          ))
        )}
      </div>
    );
  }
}