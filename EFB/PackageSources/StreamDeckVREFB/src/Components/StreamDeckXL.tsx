import {TTButton, GamepadUiView, RequiredProps, TVNode, UiViewProps } from "@efb/efb-api";
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
      <div ref={this.gamepadUiViewRef}>
        <div className="streamdeck-box">
          This is a simple box.
        </div>
      </div>
    );
  }
} 