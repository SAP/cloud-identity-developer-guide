// TEST FOR DCL FEATURES
@odm.type: 'abc'
schema {
	restr: Number,
	num1: Number,
	num2: Number,
	numArray: Number[],
	nullval: Boolean,
	id: Number,
	stringval: String,
	stringArray: String[],
	x :  {},
	bool1: Boolean,
	bool2: Boolean,
	bool3: Boolean,
	bool4: Boolean,
	boolArray: Boolean[],
	
	str1: String,
	str2: String,
	ur1: String,
	ur2: String,
	ur3: String,
	ur4: String,
	ur5: String,
	ur6: String,
	
    $same:   {
    	x: String
    },
    $Struct: {
		"Quoted-sub-name": String,
		anyOtherName: String
    },
	$Struct2: {
		"Quoted-sub-name": String,
		anyOtherName: String
    },

    "\"quoted\"": String,
    "\"quoted2\"": {
    	findme : String
    },
    
    @odm.type: 'xyz'
    @odm.hint: '123'
    element_with_annotations: String,
    
    default: String
}