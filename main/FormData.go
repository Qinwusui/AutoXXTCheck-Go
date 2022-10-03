package main

type FormData []struct {
	LinkInfo struct {
		CondFields              []interface{} `json:"condFields"`
		LinkFormType            string        `json:"linkFormType"`
		LinkFormID              int           `json:"linkFormId"`
		LinkFormValueFieldCompt string        `json:"linkFormValueFieldCompt"`
		LinkFormIDEnc           string        `json:"linkFormIdEnc"`
		LinkFormValueFieldID    int           `json:"linkFormValueFieldId"`
		Linked                  bool          `json:"linked"`
	} `json:"linkInfo"`
	Compt               string `json:"compt"`
	LoginUserForValue   bool   `json:"loginUserForValue,omitempty"`
	LayoutRatio         int    `json:"layoutRatio"`
	RelationValueConfig struct {
		CondFieldID int  `json:"condFieldId"`
		Type        int  `json:"type"`
		Open        bool `json:"open"`
	} `json:"relationValueConfig,omitempty"`
	Alias         string `json:"alias"`
	LatestValShow bool   `json:"latestValShow,omitempty"`
	OptionalScope struct {
		Options string `json:"options"`
		Type    int    `json:"type"`
	} `json:"optionalScope,omitempty"`
	ID     int `json:"id"`
	Fields []struct {
		HasDefaultValue bool `json:"hasDefaultValue"`
		Visible         bool `json:"visible"`
		Editable        bool `json:"editable"`
		Values          []struct {
			Puid  int    `json:"puid"`
			Uname string `json:"uname"`
			Val   string `json:"val"`
		} `json:"values"`
		Name   string `json:"name"`
		Verify struct {
			Required struct {
			} `json:"required"`
		} `json:"verify"`
		Tip struct {
			Imgs []interface{} `json:"imgs"`
			Text string        `json:"text"`
		} `json:"tip"`
		DefaultValueStr string `json:"defaultValueStr"`
		Label           string `json:"label"`
		SweepCode       bool   `json:"sweepCode"`
		FieldType       struct {
			Multiple bool   `json:"multiple"`
			Type     string `json:"type"`
		} `json:"fieldType"`
	} `json:"fields"`
	InDetailGroupIndex int  `json:"inDetailGroupIndex"`
	FromDetail         bool `json:"fromDetail"`
	IsShow             bool `json:"isShow"`
	HasAuthority       bool `json:"hasAuthority"`
	Formula            struct {
		SelIndex         int    `json:"selIndex"`
		CalculateFieldID string `json:"calculateFieldId"`
		Status           bool   `json:"status"`
	} `json:"formula,omitempty"`
	FormulaEdit struct {
		Formula string `json:"formula"`
	} `json:"formulaEdit,omitempty"`
	Level         int `json:"level,omitempty"`
	LocationScope struct {
		LinkedInfo struct {
		} `json:"linkedInfo"`
		MapValue     []interface{} `json:"mapValue"`
		DefaultRange int           `json:"defaultRange"`
		Select       bool          `json:"select"`
		Type         int           `json:"type"`
	} `json:"locationScope,omitempty"`
	DistanceRange          int  `json:"distanceRange,omitempty"`
	DefaultValueConfig     int  `json:"defaultValueConfig,omitempty"`
	LocationValue          int  `json:"locationValue,omitempty"`
	InDetailGroupGeneralID int  `json:"inDetailGroupGeneralId,omitempty"`
	OptionScoreShow        bool `json:"optionScoreShow,omitempty"`
	OptionScoreUsed        bool `json:"optionScoreUsed,omitempty"`
	OptionSort             struct {
		ID   string `json:"id"`
		Sort string `json:"sort"`
	} `json:"optionSort,omitempty"`
	OtherAllowed       bool `json:"otherAllowed,omitempty"`
	OptionColor        bool `json:"optionColor,omitempty"`
	OptionsLoadFromURL struct {
		IsLoadFromURL bool          `json:"isLoadFromUrl"`
		Response      []interface{} `json:"response"`
		URL           []interface{} `json:"url"`
		URLHeaders    []interface{} `json:"urlHeaders"`
	} `json:"optionsLoadFromUrl,omitempty"`
	OpenOtherOption bool `json:"openOtherOption,omitempty"`
	OptionBindInfo  struct {
		BindFormIDEnc   string `json:"bindFormIdEnc"`
		BindFieldID     int    `json:"bindFieldId"`
		BindFieldIdx    int    `json:"bindFieldIdx"`
		BindFormType    string `json:"bindFormType"`
		IsBinded        bool   `json:"isBinded"`
		BindFormID      int    `json:"bindFormId"`
		OriginalOptions []struct {
			IDArr   []interface{} `json:"idArr"`
			Score   int           `json:"score"`
			Checked bool          `json:"checked"`
			Title   string        `json:"title"`
		} `json:"originalOptions"`
		BindFieldCompt string `json:"bindFieldCompt"`
	} `json:"optionBindInfo,omitempty"`
}
