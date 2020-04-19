package pdfrender

var (
//	aquaServer = "[Aqua Server URL]"
//	imageSummary = "[This section contains the image summary: musl libc through 1.1.23 has an x87 floating-point stack adjustment imbalance, related to the math/i386/ directory. In some cases, use of this library could introduce out-of-bounds writes that are not present in an application's source code.][This section contains the image summary: musl libc through 1.1.23 has an x87 floating-point stack adjustment imbalance, related to the math/i386/ directory. In some cases, use of this library could introduce out-of-bounds writes that are not present in an application's source code.]"
	//imageName = "[img name]"
	//registry = "[image Registry]"
//	imageData = "[Image Date]"
	//imageOs = "[Image OS]"
	//imageAllowed = "[[Allowed/Disallowed]]"

//	countCritical = "[count of critical]"
//	countHigh = "[count of high]"
	//countMedium = "[count of medium]"
	//countLow = "[count of low]"
	//countNegligible = "[count of negligible]"

//	passFail = "[PASS/FAIL]"

	vulnDescription = "The application does not prevent browsers from sending sensitive information to third party sites in the referer header, despite set ting a Referrer Policy.\nWith the current Referrer Policy, when a user clicks a link that takes him to another origin (domain), the browser will add a refere r header with the URL from which he is coming from. That URL may contain sensitive information, such as password recovery toke ns or personal information, and it will be visible that other origin. For instance, if the user is at example.com/password_recovery? unique_token=14f748d89d and clicks a link to example-analytics.com, that origin will receive the complete password recovery UR L in the headers and might be able to set the users password. The same happens for requests made automatically by the applicati on, such as XHR ones.\nApplications should set a secure referrer policy that prevents sensitive data from being sent to third party sites."

	cveNumber = "[CVE-2020-1234]"
	negligible = "[Negligible]"
	cvssScrore = "[CVSS Score 3.1]"
	solution = "[Upgrade package python to version 2.6.6-68.el6_10 or above.]"

	resource = "[Python]"
	resourceFullName = "[python 2.6.6-66.el6_8]"
	fixedVersion = "[2.6.6-68.el6_10]"

	malwareListData = "[List of sensitive data][This section contains the image summary: musl libc through 1.1.23 has an x87 floating-point stack adjustment imbalance, related to the math/i386/ directory. In some cases, use of this library could introduce out-of-bounds writes that are not present in an application's source code.][This section contains the image summary: musl libc through 1.1.23 has an x87 floating-point stack adjustment imbalance, related to the math/i386/ directory. In some cases, use of this library could introduce out-of-bounds writes that are not present in an application's source code.]"
)
