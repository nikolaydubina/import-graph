package iggo_test

// GoPackageTestRunResultOutput is json formatted output result of a run go test
// {"Time":"2021-04-18T17:21:45.226785+08:00","Action":"pass","Package":"cloud.google.com/go/texttospeech/apiv1","Test":"TestTextToSpeechSynthesizeSpeechError","Elapsed":0}
//
// {"Time":"2021-04-18T17:21:45.226921+08:00","Action":"output","Package":"cloud.google.com/go/texttospeech/apiv1","Output":"coverage: 77.9% of statements\n"}
//
// {"Time":"2021-04-18T17:20:57.986807+08:00","Action":"fail","Package":"cloud.google.com/go/containeranalysis/apiv1","Test":"TestIntegration","Elapsed":0}
//
// {"Time":"2021-04-18T17:21:43.658225+08:00","Action":"output","Package":"cloud.google.com/go/talent/apiv4","Output":"testing: warning: no tests to run\n"}
//
// {"Time":"2021-04-18T17:21:45.232257+08:00","Action":"output","Package":"cloud.google.com/go/third_party/pkgsite","Output":"?   \tcloud.google.com/go/third_party/pkgsite\t[no test files]\n"}
// {"Time":"2021-04-18T17:21:45.232286+08:00","Action":"skip","Package":"cloud.google.com/go/third_party/pkgsite","Elapsed":0}
