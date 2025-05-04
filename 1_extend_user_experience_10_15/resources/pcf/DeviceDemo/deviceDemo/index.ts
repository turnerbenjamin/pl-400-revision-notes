import { IInputs, IOutputs } from "./generated/ManifestTypes";

interface PacError {
    message?: string;
    errorCode?: string;
}

enum ValueSource{
    Audio = "Audio",
    Image = "Image",
    Video = "Video",
    Barcode = "Barcode",
    Location = "Location",
    File = "File",
}

enum BaseMimeType{
    Image = "image",
    Audio = "audio",
    Video = "video",
    Unknown = "application"
}

export class deviceDemo implements ComponentFramework.StandardControl<IInputs, IOutputs> {
   
    private _notifyOutputChanged: ()=>void;
    private _deviceApi: ComponentFramework.Device; 
    private _cleanUpFunctions: (()=>void)[] = [];
    private _handlers: Record<ValueSource, ()=>Promise<string>>

    private _componentWrapper: HTMLDivElement;

    private _errorMessage = "";
    private _value = "";
    private _valueSource: string;


    /**
     * Initializes the PCF control instance.
     * 
     * @param context Property bag and component meta data
     * @param notifyOutputChanged Callback to notify the framework that the 
     *                            control's outputs have changed.
     * @param state A dictionary object that can be used to store state across 
     *              control sessions.
     * @param container The HTML element that will contain the control's UI 
     *                  elements.
     * 
     * This method sets up the initial state of the control, binds the 
     * notifyOutputChanged callback, initializes the device API, and sets up the 
     * event handlers for the control's buttons.
     */
    public init(
        context: ComponentFramework.Context<IInputs>,
        notifyOutputChanged: () => void,
        state: ComponentFramework.Dictionary,
        container: HTMLDivElement
    ): void {

        this._notifyOutputChanged = notifyOutputChanged.bind(this);
        this._deviceApi = context.device;

        this._value = "";
        this._errorMessage = "";

        this._handlers = {
            [ValueSource.Audio]: this.captureAudio.bind(this),
            [ValueSource.Image]: this.captureImage.bind(this),
            [ValueSource.Video]: this.captureVideo.bind(this),
            [ValueSource.Barcode]: this.captureBarcode.bind(this),
            [ValueSource.Location]: this.captureLocation.bind(this),
            [ValueSource.File]: this.captureFile.bind(this),
        }

        this.initialiseComponents(); 
        
        container.appendChild(this._componentWrapper);
    }


    /**
     * Updates the control view. This method is called when any value in the 
     * property bag has changed.
     * 
     * @param context - The context object provides access to the environment 
     *                  in which the control operates.
     * 
     * This method is used to update the control's UI based on the current 
     * context. In this implementation, no UI updates are required, so the 
     * method is empty.
     */
    public updateView(context: ComponentFramework.Context<IInputs>): void {
        //No UI updates required
    }


    /**
     * Returns the outputs of the control. This method is called by the 
     * framework when the control's outputs need to be retrieved.
     * 
     * @returns An object containing the current values of the control's output 
     * properties.
     */
    public getOutputs(): IOutputs {
        return {
            capturedValue: this._value,
            capturedValueSource: this._valueSource,
            errorMessage: this._errorMessage,
        };
    }


    /**
     * Cleans up the control instance. This method is called by the framework 
     * when the control is to be removed from the DOM tree. 
     *
     * Clean up functions are defined where resources are initialised 
     */
    public destroy(): void {
        for(const cleanUpFunction of this._cleanUpFunctions){
            cleanUpFunction();
        }
    }

    // Helper function to add clean up methods. Added for readability
    private defer(cleanUpCallback: ()=>void){
        this._cleanUpFunctions.push(cleanUpCallback);
    }

    // Initialise the component wrapper and appends the 6 buttons used to demo
    // device functionality
    private initialiseComponents(){
        this._componentWrapper = document.createElement("div");

        this._componentWrapper.append(
            this.buildButton.bind(this)("CAPTURE AUDIO", ValueSource.Audio),
            this.buildButton.bind(this)("CAPTURE IMAGE", ValueSource.Image),
            this.buildButton.bind(this)("CAPTURE VIDEO", ValueSource.Video),
            this.buildButton.bind(this)("CAPTURE BARCODE", ValueSource.Barcode),
            this.buildButton.bind(this)("CAPTURE LOCATION", ValueSource.Location),
            this.buildButton.bind(this)("CAPTURE FILE", ValueSource.File),
            );
    }

    // Create a build button with a click handler. The src data attribute is 
    // used by the click handler which acts as a control routing to specific
    // handlers
    private buildButton(txt: string, src: string){
        const button = document.createElement("button");
        button.textContent = txt;
        button.dataset.src = src;

        const clickHandler = this.handleButtonClick.bind(this);
        button.addEventListener("click", clickHandler);
        this.defer(()=>button.removeEventListener("click", clickHandler));

        return button;
    }

    // Control for all click handler events. Uses the src data attribute to 
    // route to the appropriate logic to extract a value
    private async handleButtonClick(e: Event){
        this._errorMessage = "";
        const target = e.target as HTMLButtonElement;
        const src = target.dataset.src as ValueSource;
        try{
        this._value = await this._handlers[ValueSource[src]]();
        }catch(error){
            this._errorMessage = this.parseErrorMessage(error as PacError);
        }
        finally{
            this._valueSource = src
            this._notifyOutputChanged();
        }
        
    }

    // Capture audio from device and return audio url. Untested as not available
    // on my phone
    private async captureAudio(){
        const audio = await this._deviceApi.captureAudio();
        return  this.getFileUrl(audio, BaseMimeType.Audio);
    }

    // Capture image from device and store image url in value field
    private async captureImage(){
        const image = await this._deviceApi.captureImage(
            { 
                height: 250, 
                width: 400, 
                allowEdit: true, 
                preferFrontCamera: false, 
                quality: 100 
            });
        return  this.getFileUrl(image, BaseMimeType.Image);
    }

    // Capture video from a device and return video url. Untested as feature is
    // not available on my phone
    private async captureVideo(){
        const video = await this._deviceApi.captureVideo();
        return this.getFileUrl(video, BaseMimeType.Video);
    }

    // Capture a barcode from device and return as a string
    private async captureBarcode(){
        const barcode = await this._deviceApi.getBarcodeValue();
        return barcode;
    }

    // Capture location and return coords as a string. Format selected to allow
    // easy parsing with PowerFx
    private async captureLocation(){
        const pos = await this._deviceApi.getCurrentPosition();
        return `${pos.coords.latitude},${pos.coords.longitude}`;
    }

    // Capture file from device and return file url
    private async captureFile(){
        const files = await this._deviceApi.pickFile({
            allowMultipleFiles: false,
            accept: "image",
            maximumAllowedFileSize: 1024 * 1024 * 10
        });
        if(!files?.length) return "";
        return  this.getFileUrl(files[0], BaseMimeType.Unknown);
    }

    

    // Helper method to create file url from a file
    private getFileUrl(file: ComponentFramework.FileObject, baseMimeType: BaseMimeType): string{
        if(!file?.fileName || !file?.fileContent){
            return "";
        };

        const fileExtension = file.fileName.split(".").pop();
        const fileContent = file.fileContent;
        const mimeType = this.getMimeType(baseMimeType, fileExtension ?? "");

        return `data:${mimeType};base64, ${fileContent}`;
    }

    // Helper method to generate mime type for file url
    private getMimeType(baseType: BaseMimeType, fileExtension: string){
        switch(baseType){
            case BaseMimeType.Image:
            case BaseMimeType.Audio:
            case BaseMimeType.Video:
                return `${baseType}/${fileExtension}`
            default:
                return `application/octet-stream`
        }

    }

    // Parse error message.
    private parseErrorMessage(error: PacError){
        return error?.message ?? error?.errorCode ?? "Error";
    }
}
