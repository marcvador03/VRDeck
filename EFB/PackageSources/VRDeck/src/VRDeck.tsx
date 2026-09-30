import {
  App,
  AppBootMode,
  AppInstallProps,
  AppSuspendMode,
  AppView,
  AppViewProps,
  Efb,
  RequiredProps,
  TVNode,
} from "@efb/efb-api";
import { FSComponent, VNode } from "@microsoft/msfs-sdk";
import { VRDeckXL } from "./Components/VRDeckXL";

import "./VRDeck.scss";

declare const BASE_URL: string;
type ConnStatus = "connecting" | "open" | "closed" | "closing";

class VRDeckAppView extends AppView<RequiredProps<AppViewProps, "bus">> {
  protected defaultView = "VRDeckXL";
  protected registerViews(): void {
    this.appViewService.registerPage("VRDeckXL", () => (
      <VRDeckXL appViewService={this.appViewService} bus={this.bus} title="StreamDeck VR"/>
    ));
    this.defaultView = "VRDeckXL";
  }

  public render(): VNode {
    return <div class="streamdeckvr-efb">{super.render()}</div>;
  }

  private socket: WebSocket | null = null;
  private port: number = 8081;
  private cleanClose:boolean = false;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;

  private handlerOpen = (): void => {  
    console.log("[StreamDeckVR] WebSocket connection established")
    this.publishStatus();
  };

  private handlerError = (event: Event): void => {  
    console.error(`[StreamDeckVR] WebSocket error: ${event}`)
  };

  private handlerClose = (event: CloseEvent): void => {  
    console.log(`[StreamDeckVR] connection closed with code: ${event.code} and reason: ${event.reason}`)
    this.socket = null;
    if (this.cleanClose == false) {
      this.reconnectTimer = setTimeout(() => this.connect(), 2000);
    }
    this.publishStatus();
  };

  private handlerMsg = (event: MessageEvent): void => {  
    console.log(`[StreamDeckVR] Raw message received: ${event.data}`);
      try {
        const msgJson = JSON.parse(event.data);
        console.log(`[StreamDeckVR] JSon parsed ${msgJson}`);
        this.bus.pub("streamdeck-labels-updated", msgJson);
        this.appViewService.update(performance.now());
      } catch (err) {
        console.error("[StreamDeckVR] Failed to parse JSON:", err);
        console.error("[StreamDeckVR] Raw data that failed to parse:", event.data);
      };
  };

  public async onOpen(): Promise<void> {
    this.connect();
    this.publishStatus();
  }

  public async onClose(): Promise<void> {
    this.cleanClose = true;
    this.publishStatus();
    if (this.reconnectTimer != null) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    };
    this.closesocket();
  }
  
  public connect(): void {
    if (this.socket != null) {
      console.error("[StreamDeckVR] Cannot initiate WebSocket connection, existing connecting still live")
      return;
    }
    this.cleanClose = false;
    console.log("[StreamDeckVR] Attempting to connect to WebSocket at ws://localhost:8081/streamdeckvr...");
    try {
      this.socket = new WebSocket(`ws://localhost:${this.port}/streamdeckvr`)
      this.socket.onopen = this.handlerOpen;
      this.socket.onclose = this.handlerClose;
      this.socket.onerror = this.handlerError;
      this.socket.onmessage = this.handlerMsg;
    } catch (err) {
      console.error(`Error while initializing WebSocket at: ws://localhost:${this.port}/streamdeckvr with error ${err}`);
    };
    this.publishStatus();
  };

  public closesocket(): void {
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    };
  };

  public publishStatus(): void {
    const state = this.getStatus();
    this.bus.pub("streamdeck-connection-status", this.getStatus());
   

  };

  private getStatus(): ConnStatus {
    if (this.socket == null)
        return "closed";
    else {
      if (this.socket.readyState == WebSocket.OPEN)
        return "open";
      else if (this.socket.readyState == WebSocket.CLOSED)
        return "closed";
      else if (this.socket.readyState == WebSocket.CONNECTING)
        return "connecting";
      else if (this.socket.readyState ==  WebSocket.CLOSING)
        return "closing";
    }
    return "closed";
  }

  //   public async onPause(): Promise<void> {
  //   if (this.socket) {
  //     this.socket.close();
  //     this.socket = null;
  //     console.log("disconnected")
  //   }
  // }
}

class VRDeck extends App {
 
  public get name(): string {
    return VRDeck.name;
  }
 
  public get icon(): string {
    return `${BASE_URL}/Assets/app-icon.svg`;
  }

  /**
   * Optional attribute
   * Allow to choose BootMode between COLD / WARM / HOT
   * Default behavior : AppBootMode.COLD
   *
   * COLD : No dom preloaded in memory
   * WARM : App -> AppView are loaded but not rendered into DOM
   * HOT : App -> AppView -> Pages are rendered and injected into DOM
   */
  public BootMode = AppBootMode.COLD;

  /**
   * Optional attribute
   * Allow to choose SuspendMode between SLEEP / TERMINATE
   * Default behavior : AppSuspendMode.SLEEP
   *
   * SLEEP : Default behavior, does nothing, only hiding and sleeping the app if switching to another one
   * TERMINATE : Hiding the app, then killing it by removing it from DOM (BootMode is checked on next frame to reload it and/or to inject it, see BootMode)
   */
  public SuspendMode = AppSuspendMode.SLEEP;

  /**
   * Optional method
   * Allow to resolve some dependencies, install external data, check an api key, etc...
   * @param _props props used when app has been setted up.
   * @returns Promise<void>
   */
  public async install(_props: AppInstallProps): Promise<void> { 
    Efb.loadCss(`${BASE_URL}/VRDeck.css`);
    return Promise.resolve();
  }

  public render(): TVNode<VRDeckAppView> {
    return <VRDeckAppView bus={this.bus} />;
  }
}

Efb.use(VRDeck);
