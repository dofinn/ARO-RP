package e2e

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	conditions "github.com/serge1peshcoff/selenium-go-conditions"
	"github.com/tebeka/selenium"
)

var _ = Describe("Admin Portal E2E Testing", func() {
	BeforeEach(skipIfNotInDevelopmentEnv)
	var wdPoint *selenium.WebDriver
	var wd selenium.WebDriver
	var host string

	JustBeforeEach(func() {
		host, wdPoint = adminPortalSessionSetup()
		wd = *wdPoint
		wd.Get(host + "/v2")
		wd.Refresh()
	})

	AfterEach(func() {
		if wd != nil {
			wd.Quit()
		}
	})

	It("Should be able to populate cluster data correctly", func() {
		err := wd.Wait(conditions.ElementIsLocated(selenium.ByCSSSelector, "div[data-automation-key='name']"))
		if err != nil {
			SaveScreenshotAndExit(wd, err)
		}

		cluster, err := wd.FindElement(selenium.ByCSSSelector, "div[data-automation-key='name']")
		if err != nil {
			SaveScreenshotAndExit(wd, err)
		}

		Expect(cluster.Text()).To(Equal(os.Getenv("CLUSTER")))
	})

	It("Should be able to filter cluster data correctly", func() {
		wd.Wait(conditions.ElementIsLocated(selenium.ByCSSSelector, "div[data-automation-key='name']"))

		filter, err := wd.FindElement(selenium.ByCSSSelector, "input[placeholder='Filter on resource ID']")
		if err != nil {
			SaveScreenshotAndExit(wd, err)
		}

		// Set filter so it doesn't match cluster name
		filter.SendKeys("Incorrect Cluster")

		wd.Wait(conditions.ElementIsLocated(selenium.ByID, "ClusterCount"))
		text, err := wd.FindElement(selenium.ByID, "ClusterCount")
		if err != nil {
			SaveScreenshotAndExit(wd, err)
		}

		Expect(text.Text()).To(Equal("Showing 0 items"))
	})

	It("Should be able to populate cluster info panel correctly", func() {
		const CLUSTER_INFO_HEADINGS = 17

		err := wd.Wait(conditions.ElementIsLocated(selenium.ByCSSSelector, "div[data-automation-key='name']"))
		if err != nil {
			SaveScreenshotAndExit(wd, err)
		}

		cluster, err := wd.FindElement(selenium.ByCSSSelector, "div[data-automation-key='name']")
		if err != nil {
			SaveScreenshotAndExit(wd, err)
		}

		err = cluster.Click()
		if err != nil {
			SaveScreenshotAndExit(wd, err)
		}

		err = wd.WaitWithTimeout(conditions.ElementIsLocated(selenium.ByID, "ClusterDetailCell"), 2*time.Minute)
		if err != nil {
			SaveScreenshotAndExit(wd, err)
		}

		panelSpans, err := wd.FindElements(selenium.ByID, "ClusterDetailCell")
		if err != nil {
			SaveScreenshotAndExit(wd, err)
		}

		Expect(len(panelSpans)).To(Equal(CLUSTER_INFO_HEADINGS * 3))

		panelFields := panelSpans[0 : CLUSTER_INFO_HEADINGS-1]
		panelColons := panelSpans[CLUSTER_INFO_HEADINGS : CLUSTER_INFO_HEADINGS*2-1]
		panelValues := panelSpans[CLUSTER_INFO_HEADINGS*2 : len(panelSpans)-1]

		for _, panelField := range panelFields {
			panelText, err := panelField.Text()
			if err != nil {
				SaveScreenshotAndExit(wd, err)
			}

			Expect(panelText).To(Not(Equal("")))
		}

		for _, panelField := range panelColons {
			panelText, err := panelField.Text()
			if err != nil {
				SaveScreenshotAndExit(wd, err)
			}

			Expect(panelText).To(Equal(":"))
		}

		for _, panelField := range panelValues {
			panelText, err := panelField.Text()
			if err != nil {
				SaveScreenshotAndExit(wd, err)
			}

			Expect(panelText).To(Not(Equal("")))
		}
	})

	It("Should be able to copy cluster resource id", func() {
		wd.Wait(conditions.ElementIsLocated(selenium.ByCSSSelector, "button[aria-label='Copy Resource ID']"))

		button, err := wd.FindElement(selenium.ByCSSSelector, "button[aria-label='Copy Resource ID']")
		if err != nil {
			SaveScreenshotAndExit(wd, err)
		}

		button.Click()

		filter, err := wd.FindElement(selenium.ByCSSSelector, "input[placeholder='Filter on resource ID']")
		if err != nil {
			SaveScreenshotAndExit(wd, err)
		}

		// Paste clipboard
		filter.Click()
		filter.SendKeys(selenium.ControlKey + "v")
		resourceId, err := filter.GetAttribute("value")

		if err != nil {
			SaveScreenshotAndExit(wd, err)
		}

		Expect(resourceId).To(ContainSubstring("/providers/Microsoft.RedHatOpenShift/openShiftClusters/" + os.Getenv("CLUSTER")))
	})

	It("Should be able to open ssh panel and get ssh details", func() {
		wd.Wait(conditions.ElementIsLocated(selenium.ByCSSSelector, "button[aria-label='SSH']"))

		button, err := wd.FindElement(selenium.ByCSSSelector, "button[aria-label='SSH']")
		if err != nil {
			SaveScreenshotAndExit(wd, err)
		}

		button.Click()

		wd.Wait(conditions.ElementIsLocated(selenium.ByID, "sshModal"))
		wd.Wait(conditions.ElementIsLocated(selenium.ByID, "sshDropdown"))

		sshDropdown, err := wd.FindElement(selenium.ByID, "sshDropdown")
		if err != nil {
			SaveScreenshotAndExit(wd, err)
		}

		sshDropdown.Click()

		wd.Wait(conditions.ElementIsLocated(selenium.ByID, "sshDropdown-list0"))
		machine, err := wd.FindElement(selenium.ByID, "sshDropdown-list0")
		if err != nil {
			SaveScreenshotAndExit(wd, err)
		}

		machine.Click()

		wd.Wait(conditions.ElementIsLocated(selenium.ByID, "sshButton"))
		requestBtn, err := wd.FindElement(selenium.ByID, "sshButton")
		if err != nil {
			SaveScreenshotAndExit(wd, err)
		}

		requestBtn.Click()

		// Test fails if these labels aren't present
		err = wd.Wait(conditions.ElementIsLocated(selenium.ByID, "sshCommand"))
		if err != nil {
			SaveScreenshotAndExit(wd, err)
		}
	})
})