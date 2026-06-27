import {TTButton, GamepadUiView, RequiredProps, TVNode, UiViewProps } from "@efb/efb-api";
import { FSComponent } from "@microsoft/msfs-sdk";
import "./StreamDeckXL.scss";

interface StreamDeckXLProps extends RequiredProps<UiViewProps, "appViewService"> {
  title?: string;
}

export class StreamDeckXL extends GamepadUiView<HTMLDivElement, StreamDeckXLProps> {
  public readonly tabName = StreamDeckXL.name;

  public render(): TVNode<HTMLDivElement> {
    return (
      <div ref={this.gamepadUiViewRef} className="streamdeck-container">
         <div class="header">
          <TTButton
            key="Go back"
            type="secondary"
            callback={(): void => {
              this.props.appViewService.goBack();
            }}
          />
          <h2>{this.props.title}</h2>
        </div>
        
      </div>
    );
  }
}