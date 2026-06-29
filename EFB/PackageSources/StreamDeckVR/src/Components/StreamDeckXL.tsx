import { GamepadUiView, RequiredProps, TVNode, UiViewProps } from "@efb/efb-api";
import { FSComponent } from "@microsoft/msfs-sdk";
import "./StreamDeckXL.scss";

interface StreamDeckXLProps extends RequiredProps<UiViewProps, "appViewService"> {
  title?: string;
  color?: string;
  labelsData: { buttons: { row: number; col: number; label: string }[] };
}

export class StreamDeckXL extends GamepadUiView<HTMLDivElement, StreamDeckXLProps> {
  public readonly tabName = StreamDeckXL.name;
  private cellLabels: string[][] = Array.from({ length: 4 }, (_, row) =>
    Array.from({ length: 8 }, (_, col) => String(row * 8 + col + 1))
  );
  private previousLabelsData: { buttons: { row: number; col: number; label: string }[] } | null = null;

  public onUpdate(): void {
    if (!this.props.labelsData) return;

    const isDifferent = !this.previousLabelsData ||
      JSON.stringify(this.previousLabelsData) !== JSON.stringify(this.props.labelsData);

    if (isDifferent) {
      this.updateLabels(this.props.labelsData);
      this.previousLabelsData = this.props.labelsData;
    }
  }

  public updateLabels(json: { buttons: { row: number; col: number; label: string }[] }): void {
    console.log("Update triggered", this.props.labelsData)
    this.cellLabels = Array.from({ length: 4 }, (_, row) =>
      Array.from({ length: 8 }, (_, col) => String(row * 8 + col + 1))
    );
    json.buttons.forEach(button => {
      if (button.row >= 0 && button.row < 4 && button.col >= 0 && button.col < 8) {
        this.cellLabels[button.row][button.col] = button.label;
      }
    });
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