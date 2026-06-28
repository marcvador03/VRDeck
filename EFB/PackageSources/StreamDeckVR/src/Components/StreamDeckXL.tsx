import { GamepadUiView, RequiredProps, TVNode, UiViewProps } from "@efb/efb-api";
import { FSComponent } from "@microsoft/msfs-sdk";
import "./StreamDeckXL.scss";

interface StreamDeckXLProps extends RequiredProps<UiViewProps, "appViewService"> {
  title?: string;
  color?: string;
}

export class StreamDeckXL extends GamepadUiView<HTMLDivElement, StreamDeckXLProps> {
  public readonly tabName = StreamDeckXL.name;

  public render(): TVNode<HTMLDivElement> {
    return (
      <div ref={this.gamepadUiViewRef} class="streamdeck-container">
        {Array.from({ length: 32 }).map((_, index) => (
          <div key={`cell-${index}`} class="streamdeck-cell">
            {index + 1} {/* Button label (1-32) */}
          </div>
        ))}
      </div>
    );
  }
}